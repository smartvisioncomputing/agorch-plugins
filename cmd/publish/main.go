// Command publish signs the sealed packages in index/ and publishes them to the
// Agorch registry through the ordinary publisher path (account, namespace claim,
// key, artifact, listing, release, checks, review, stable channel). It drives
// the API from the state it reads, so a second run changes nothing.
//
// Signing reuses ao-packages/catalog/build (the same Sign the registry's own
// first-party tests use). Credentials come from the environment and are never
// printed: AGORCH_API, AGORCH_PUBLISHER_TOKEN, AGORCH_PLATFORM_TOKEN,
// AGORCH_PUBLISHER_SEED_FILE, AGORCH_OWNER_EMAIL (see scripts/publish.sh).
package main

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/smartvisioncomputing/aicl/pkg/artifacttrust/signer"
	trustv1 "github.com/smartvisioncomputing/aicl/pkg/artifacttrust/v1"
	componentv11 "github.com/smartvisioncomputing/aicl/pkg/componentpackage/v11"
	catalogbuild "github.com/smartvisioncomputing/ao-packages/catalog/build"
)

type obj = map[string]any

type reply struct {
	status   int
	body     []byte
	replayed bool
}

func (r reply) code() string { // the API's machine-readable error code, if any
	var e struct {
		Code string `json:"code"`
	}
	_ = json.Unmarshal(r.body, &e)
	return e.Code
}

func (r reply) obj() obj {
	var o obj
	_ = json.Unmarshal(r.body, &o)
	return o
}

func str(o obj, k string) string { s, _ := o[k].(string); return s }

// items is the list a collection reply carries; a reply without one is an empty list, never a panic.
func items(r reply) []any { l, _ := r.obj()["items"].([]any); return l }

type publisher struct {
	base, publisherTok, platformTok string
	pid, ns, owner, keyID           string
	signer                          signer.Signer
	changes                         int
	pending, blocked                []string
}

func (p *publisher) do(method, path, tok, key, ifMatch string, body any) (reply, error) {
	var rd io.Reader
	if body != nil {
		var b bytes.Buffer
		enc := json.NewEncoder(&b)
		enc.SetEscapeHTML(false) // the registry demands the canonical bytes: "<" must stay "<"
		if err := enc.Encode(body); err != nil {
			return reply{}, err
		}
		rd = &b
	}
	req, err := http.NewRequest(method, p.base+path, rd)
	if err != nil {
		return reply{}, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	if key != "" {
		req.Header.Set("Idempotency-Key", key)
	}
	if ifMatch != "" {
		req.Header.Set("If-Match", ifMatch)
	}
	res, err := (&http.Client{Timeout: 60 * time.Second}).Do(req)
	if err != nil {
		return reply{}, err
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	return reply{res.StatusCode, raw, res.Header.Get("Idempotency-Replayed") == "true"}, nil
}

// must runs a call that has to succeed with one of the given statuses.
func (p *publisher) must(method, path, tok, key, ifMatch string, body any, want ...int) (reply, error) {
	r, err := p.do(method, path, tok, key, ifMatch, body)
	if err != nil {
		return r, err
	}
	for _, w := range want {
		if r.status == w {
			return r, nil
		}
	}
	return r, fmt.Errorf("%s %s: HTTP %d %s", method, path, r.status, r.code())
}

func (p *publisher) logf(format string, a ...any) { fmt.Printf(format+"\n", a...) }

// created records a mutation that really changed state (a replay does not count).
func (p *publisher) created(r reply, what string) {
	if r.replayed {
		p.logf("  = %s (replay)", what)
		return
	}
	p.changes++
	p.logf("  + %s", what)
}

func (p *publisher) bootstrap() error {
	acct := "/publisher-accounts/" + p.pid
	r, err := p.do("GET", acct, p.publisherTok, "", "", nil)
	if err != nil {
		return err
	}
	if r.status == http.StatusNotFound {
		if p.owner == "" {
			return errors.New("AGORCH_OWNER_EMAIL is required to create the publisher account")
		}
		c, err := p.must("POST", "/platform/publisher-accounts", p.platformTok, "seed-account-"+p.pid, "", obj{
			"publisher_id": p.pid, "display_name": "Agorch", "publisher_type": "first_party",
			"support_email": "contact@agorch.ai", "security_contact": "contact@agorch.ai",
			"homepage_url": "https://agorch.ai", "initial_owner": p.owner}, 201)
		if err != nil {
			return err
		}
		p.created(c, "publisher account "+p.pid)
		r, err = p.must("GET", acct, p.publisherTok, "", "", nil, 200)
		if err != nil {
			return err
		}
	} else if r.status != 200 {
		return fmt.Errorf("GET %s: HTTP %d %s", acct, r.status, r.code())
	}
	if a := r.obj(); str(a, "status") != "active" {
		c, err := p.must("POST", "/platform/"+strings.TrimPrefix(acct, "/")+"/actions/activate", p.platformTok, "seed-activate-"+p.pid, str(a, "etag"), obj{"reason": "first-party seed"}, 200)
		if err != nil {
			return err
		}
		p.created(c, "publisher account activated")
	}

	// namespace claim
	r, err = p.must("GET", acct+"/namespace-claims", p.publisherTok, "", "", nil, 200)
	if err != nil {
		return err
	}
	var claim obj
	for _, it := range items(r) {
		if o := it.(obj); str(o, "namespace") == p.ns {
			claim = o
		}
	}
	if claim == nil {
		c, err := p.must("POST", acct+"/namespace-claims", p.publisherTok, "seed-claim-"+p.ns, "", obj{"namespace": p.ns, "claim_method": "first_party", "claim_subject": p.pid}, 201)
		if err != nil {
			return fmt.Errorf("%w (the API must admit a first_party claim for %s/%s)", err, p.pid, p.ns)
		}
		p.created(c, "namespace claim "+p.ns)
		claim = c.obj()
	}
	if str(claim, "status") != "verified" {
		c, err := p.must("POST", "/platform"+acct+"/namespace-claims/"+p.ns+"/actions/verify", p.platformTok, "seed-verify-"+p.ns+"-"+str(claim, "etag"), str(claim, "etag"), obj{"reason": "first-party seed"}, 200)
		if err != nil {
			return err
		}
		p.created(c, "namespace claim verified")
	}

	// publisher signing key
	pub, err := p.signer.PublicKey(context.Background())
	if err != nil {
		return err
	}
	r, err = p.must("GET", acct+"/keys", p.publisherTok, "", "", nil, 200)
	if err != nil {
		return err
	}
	have := false
	for _, it := range items(r) {
		if o := it.(obj); str(o, "key_id") == p.keyID && str(o, "status") == "active" {
			have = true
		}
	}
	if !have {
		c, err := p.must("POST", acct+"/keys", p.publisherTok, "seed-key-"+p.keyID, "", obj{
			"namespace": p.ns, "key_id": p.keyID, "algorithm": trustv1.AlgorithmEd25519,
			"public_key_base64": signer.PublicKeyBase64(pub),
			"not_before":        time.Now().Add(-time.Minute).UTC().Format(time.RFC3339), "not_after": time.Now().AddDate(2, 0, 0).UTC().Format(time.RFC3339)}, 201)
		if err != nil {
			return err
		}
		p.created(c, "publisher key "+p.keyID)
	}
	return nil
}

// platformKeys registers the public halves of the API's own check and review
// signers, but only when they are given as local://<key-id>:<hex seed> (the form
// the API reads from AO_CATALOG_CHECK_SIGNER_KEY and AO_CATALOG_REVIEW_SIGNER_KEY).
// A kms:// signer is the operator's to register: its public key is not ours to read.
func (p *publisher) platformKeys() error {
	specs := []struct{ env, ns, signerID, role string }{
		{"AO_CATALOG_CHECK_SIGNER_KEY", "ao.trust.automated-checks", "platform:ao-automated-checks", "automated_checks"},
		{"AO_CATALOG_REVIEW_SIGNER_KEY", "ao.trust.release-review", "platform:ao-release-review", "ao_release_reviewer"},
	}
	for _, s := range specs {
		id, seedHex, ok := strings.Cut(strings.TrimPrefix(os.Getenv(s.env), "local://"), ":")
		if !ok || !strings.HasPrefix(os.Getenv(s.env), "local://") {
			p.logf("  . %s is not local://: register that platform key by hand", s.env)
			continue
		}
		seed, err := hex.DecodeString(seedHex)
		if err != nil || len(seed) != ed25519.SeedSize {
			return fmt.Errorf("%s: seed is not 64 hex characters", s.env)
		}
		r, err := p.must("GET", "/platform/catalog/trusted-keys?limit=100", p.platformTok, "", "", nil, 200)
		if err != nil {
			return err
		}
		found := false
		for _, it := range items(r) {
			if o := it.(obj); str(o, "key_id") == id && str(o, "signer_namespace") == s.ns && o["revoked_at"] == nil {
				found = true
			}
		}
		if found {
			continue
		}
		pub := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
		now := time.Now().UTC()
		c, err := p.must("POST", "/platform/catalog/trusted-keys", p.platformTok, "seed-platform-key-"+id, "", obj{
			"signer_namespace": s.ns, "signer_id": s.signerID, "role": s.role, "key_id": id, "algorithm": trustv1.AlgorithmEd25519,
			"public_key_base64": signer.PublicKeyBase64(pub),
			"not_before":        now.Add(-time.Minute).Format(time.RFC3339), "not_after": now.AddDate(2, 0, 0).Format(time.RFC3339)}, 201)
		if err != nil {
			return err
		}
		p.created(c, "platform key "+id)
	}
	return nil
}

type pkgInfo struct{ ns, id, version, digest string }

func (p *publisher) publish(file string, present obj) error {
	canonical, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	var raw struct {
		Digest  string `json:"canonical_digest"`
		Package struct {
			Namespace string `json:"namespace"`
			PackageID string `json:"package_id"`
			Version   string `json:"version"`
		} `json:"package"`
	}
	if err := json.Unmarshal(canonical, &raw); err != nil {
		return fmt.Errorf("%s: %w", file, err)
	}
	k := pkgInfo{raw.Package.Namespace, raw.Package.PackageID, raw.Package.Version, raw.Digest}
	ref := k.ns + "/" + k.id
	p.logf("%s@%s %s", ref, k.version, k.digest)
	acct := "/publisher-accounts/" + p.pid
	artifact := obj{"namespace": k.ns, "package_id": k.id, "version": k.version, "digest": k.digest}

	// 1. artifact: sign and ingest unless the registry already holds this exact digest.
	r, err := p.do("GET", fmt.Sprintf("%s/artifacts/%s/%s/versions/%s/digests/%s", acct, k.ns, k.id, k.version, k.digest), p.publisherTok, "", "", nil)
	if err != nil {
		return err
	}
	if r.status == http.StatusNotFound {
		art := catalogbuild.Artifact{Canonical: canonical, Subject: trustv1.ArtifactSubject{
			Namespace: k.ns, PackageID: k.id, Version: componentv11.SemVer(k.version),
			Digest: componentv11.Digest(k.digest), MediaType: componentv11.ContentMediaType}}
		bundle, err := catalogbuild.Sign(context.Background(), art, p.signer)
		if err != nil {
			return err
		}
		c, err := p.must("POST", acct+"/artifacts", p.publisherTok, "seed-ingest-"+k.digest, "", obj{"artifact": json.RawMessage(canonical), "signature": bundle}, 201)
		if err != nil {
			return err
		}
		p.created(c, "artifact ingested")
	} else if r.status != 200 {
		return fmt.Errorf("GET artifact: HTTP %d %s", r.status, r.code())
	}

	// 2. listing: the genesis body is constant (only the package identity), so its
	// idempotency key always replays to the same listing id; the presentation is
	// then brought up to catalog.json by PATCH, only where it differs.
	c, err := p.must("POST", acct+"/listings", p.publisherTok, "seed-listing-v1-"+ref, "", obj{
		"namespace": k.ns, "package_id": k.id, "display_name": k.id, "summary": "Agorch first-party connector.",
		"support_url": "https://agorch.ai", "documentation_url": "https://agorch.ai",
		"license_spdx": "Apache-2.0", "hosting": "publisher_hosted"}, 201)
	if err != nil {
		return err
	}
	if !c.replayed {
		p.created(c, "listing created")
	}
	lid := str(c.obj(), "listing_id")
	lpath := acct + "/listings/" + lid
	listing, err := p.must("GET", lpath, p.publisherTok, "", "", nil, 200)
	if err != nil {
		return err
	}
	if want := (obj{"display_name": present["display_name"], "summary": present["summary"]}); !sameFields(listing.obj(), want) {
		if _, err := p.patch(lpath, listing.obj(), want); err != nil {
			return err
		}
		p.changes++
		p.logf("  + listing text updated")
	}

	// 3. release -> submit -> checks -> review -> stable, each only if still due.
	rel, err := p.must("POST", lpath+"/releases", p.publisherTok, "seed-release-"+k.digest, "", obj{"artifact_ref": artifact, "requested_channel": "stable"}, 201)
	if err != nil {
		return err
	}
	if !rel.replayed {
		p.created(rel, "release created")
	}
	rpath := lpath + "/releases/" + k.digest
	cur, err := p.must("GET", rpath, p.publisherTok, "", "", nil, 200)
	if err != nil {
		return err
	}
	review := str(cur.obj(), "review_id")
	if str(cur.obj(), "review_state") == "draft" {
		s, err := p.must("POST", rpath+"/actions/submit", p.publisherTok, "seed-submit-"+k.digest, str(cur.obj(), "etag"), obj{"policy_version": "ao.catalog.publication.v1"}, 200)
		if err != nil {
			return err
		}
		p.created(s, "release submitted")
		cur = s
	}
	if str(cur.obj(), "review_state") == "submitted" {
		pass, err := p.must("POST", "/platform/catalog/check-passes", p.platformTok, "", "", obj{"limit": 20}, 200)
		if err != nil {
			return err
		}
		if n, _ := pass.obj()["completed"].(float64); n > 0 {
			p.changes++
			p.logf("  + automated checks ran")
		}
		if cur, err = p.must("GET", rpath, p.publisherTok, "", "", nil, 200); err != nil {
			return err
		}
	}
	run, _ := cur.obj()["latest_check_run"].(obj)
	if str(cur.obj(), "review_state") == "changes_requested" {
		// The checks failed: nothing more to do until the package changes. Report the codes.
		var failed []string
		results, _ := run["results"].([]any)
		for _, it := range results {
			if o, _ := it.(obj); str(o, "outcome") == "failed" {
				failed = append(failed, str(o, "code"))
			}
		}
		p.blocked = append(p.blocked, ref)
		p.logf("  ! %s blocked by automated checks: %s", ref, strings.Join(failed, ", "))
		return nil
	}
	if str(cur.obj(), "review_state") == "checking" { // checks passed, waiting for the platform review
		st, err := p.must("GET", "/platform/catalog/listings/"+lid+"/channels/stable", p.platformTok, "", "", nil, 200)
		if err != nil {
			return err
		}
		d, err := p.must("POST", "/platform/catalog/reviews/"+review+"/decisions", p.platformTok, "seed-approve-"+k.digest, str(cur.obj(), "etag"), obj{
			"decision": "approved", "reason_code": "review.approved", "explanation": "First-party package verified through the normal path.",
			"policy_version": "ao.catalog.publication.v1", "check_run_id": str(run, "check_run_id"),
			"automated_attestation_bundle_digest": str(cur.obj(), "automated_attestation_bundle_digest"),
			"expected_channel_etag":               str(st.obj(), "etag")}, 200)
		if err != nil {
			return err
		}
		p.created(d, "release approved, stable channel updated")
		cur = d
	}
	if s := str(cur.obj(), "review_state"); s != "approved" {
		return fmt.Errorf("%s: release is %q, expected approved", ref, s)
	}

	// 4. presentation fields owned by the catalog (category, rank, vendor_official,
	// beta): tolerant. An API that does not know them yet is reported, not fatal.
	return p.presentation(lpath, ref, obj{
		"category": present["category"], "rank": present["rank"],
		"vendor_official": present["vendor_official"], "beta": present["beta"]})
}

func sameFields(have, want obj) bool {
	for k, v := range want {
		a, _ := json.Marshal(have[k])
		b, _ := json.Marshal(v)
		if !bytes.Equal(a, b) {
			return false
		}
	}
	return true
}

func (p *publisher) patch(lpath string, listing, body obj) (reply, error) {
	sum := sha256.Sum256([]byte(fmt.Sprint(body)))
	return p.must("PATCH", lpath, p.publisherTok, "seed-patch-"+lpath+"-"+hex.EncodeToString(sum[:4]), str(listing, "etag"), body, 200)
}

func (p *publisher) presentation(lpath, ref string, want obj) error {
	listing, err := p.must("GET", lpath, p.publisherTok, "", "", nil, 200)
	if err != nil {
		return err
	}
	if sameFields(listing.obj(), want) {
		return nil
	}
	// One PATCH with the publisher token. The API wants a principal that is
	// a publisher member AND a platform admin for vendor_official, rank and beta.
	sum := sha256.Sum256([]byte(fmt.Sprint(want)))
	last, err := p.do("PATCH", lpath, p.publisherTok, "seed-present-"+lpath+"-"+hex.EncodeToString(sum[:4]), str(listing.obj(), "etag"), want)
	if err != nil {
		return err
	}
	if last.status == 200 {
		p.created(last, "presentation fields set")
		return nil
	}
	p.pending = append(p.pending, ref)
	p.logf("  ! presentation fields not applied (HTTP %d %s): the API does not accept them yet, or the publisher principal is not a platform admin", last.status, last.code())
	return nil
}

func main() {
	indexDir := flag.String("index", "index", "directory with the sealed packages")
	catalogFile := flag.String("catalog", "catalog.json", "catalog.json with the presentation fields")
	flag.Parse()
	if err := run(*indexDir, *catalogFile); err != nil {
		fmt.Fprintln(os.Stderr, "publish:", err)
		os.Exit(1)
	}
}

func run(indexDir, catalogFile string) error {
	need := func(k string) string {
		v := strings.TrimSpace(os.Getenv(k))
		if v == "" {
			fmt.Fprintln(os.Stderr, "publish: missing", k)
			os.Exit(2)
		}
		return v
	}
	seedHex, err := os.ReadFile(need("AGORCH_PUBLISHER_SEED_FILE"))
	if err != nil {
		return err
	}
	seed, err := hex.DecodeString(strings.TrimSpace(string(seedHex)))
	if err != nil || len(seed) != ed25519.SeedSize {
		return errors.New("AGORCH_PUBLISHER_SEED_FILE must hold 64 hex characters")
	}
	var cat struct {
		Namespace   string `json:"namespace"`
		PublisherID string `json:"publisher_id"`
		Connectors  []obj  `json:"connectors"`
	}
	raw, err := os.ReadFile(catalogFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, &cat); err != nil {
		return err
	}
	p := &publisher{
		base: strings.TrimRight(need("AGORCH_API"), "/"), publisherTok: need("AGORCH_PUBLISHER_TOKEN"), platformTok: need("AGORCH_PLATFORM_TOKEN"),
		pid: cat.PublisherID, ns: cat.Namespace, owner: os.Getenv("AGORCH_OWNER_EMAIL"), keyID: os.Getenv("AGORCH_PUBLISHER_KEY_ID"),
	}
	if p.keyID == "" {
		p.keyID = "agorch-publisher-1"
	}
	// ponytail: local seed signer only; the production signer is a KMS key and plugs in here.
	if p.signer, err = signer.NewLocalSigner(p.keyID, seed); err != nil {
		return err
	}
	files, _ := filepath.Glob(filepath.Join(indexDir, "*.json"))
	sort.Strings(files)
	if len(files) == 0 {
		return fmt.Errorf("no packages in %s", indexDir)
	}
	byID := map[string]obj{}
	for _, c := range cat.Connectors {
		byID[str(c, "id")] = c
	}

	p.logf("bootstrap %s / %s", p.pid, p.ns)
	if err := p.bootstrap(); err != nil {
		return err
	}
	if err := p.platformKeys(); err != nil {
		return err
	}
	for _, f := range files {
		var id struct {
			Package struct {
				Namespace string `json:"namespace"`
				PackageID string `json:"package_id"`
			} `json:"package"`
		}
		b, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		_ = json.Unmarshal(b, &id)
		present, ok := byID[id.Package.Namespace+"/"+id.Package.PackageID]
		if !ok {
			return fmt.Errorf("%s is not in %s (run scripts/check.sh)", f, catalogFile)
		}
		if err := p.publish(f, present); err != nil {
			return err
		}
	}
	if p.changes == 0 {
		p.logf("done: no changes")
	} else {
		p.logf("done: %d changes", p.changes)
	}
	if len(p.pending) > 0 {
		p.logf("presentation fields pending on the API for %d packages", len(p.pending))
	}
	if len(p.blocked) > 0 {
		return fmt.Errorf("%d of %d packages are blocked by the automated checks: %s", len(p.blocked), len(files), strings.Join(p.blocked, ", "))
	}
	return nil
}
