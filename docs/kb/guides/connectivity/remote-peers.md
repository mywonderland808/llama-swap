---
title: Route to remote peers
summary: Expose models from another host, and keep X-Local-Only requests on local targets.
category: guides
tags: [peers, remote, networking, local-only]
config_keys: [peers, peers.*.proxy, peers.*.apiKey, peers.*.models, peers.*.trust]
updated: 2026-10-30
---

# Route to remote peers

Declare a peer with the other llama-swap `proxy` URL and its model list. Its models are addressed as
`peer-name/model-id`, so a peer named `sippy` serving `gemma-4-12B` appears as
`sippy/gemma-4-12B`.

```yaml
peers:
  sippy:
    proxy: http://sippy:8080
    apiKey: ${env.SIPPY_API_KEY}
    trust: local
    models: [gemma-4-12B]
  openrouter:
    proxy: https://openrouter.ai/api
    apiKey: ${env.OPENROUTER_API_KEY}
    models: [openai/gpt-4o]
```

The peer must be reachable from this host and its API key must match the remote
server. See `guides/api-integration/api-keys-and-auth` for keeping it out of
the committed file.

## trust

`trust` is `local` or `remote`. Only `local` (any case, surrounding space
ignored) is a local target. Omit the key, or set anything else, and the peer
is remote. A typo such as `locla` is remote, so it cannot receive a flagged
request. Models this instance runs itself are always local and have no `trust`
field.

Use `local` for another llama-swap you run. Leave a hosted API as remote.
Marking a hosted API `local` is an operator assertion: llama-swap will send
flagged requests there.

## X-Local-Only

Send `X-Local-Only` to refuse remote targets. The header is on when it is
present, including an empty value, `1`, or `true`. `0` and `false` (any case,
surrounding space ignored) turn it off. Any other value, such as `yes` or
`maybe`, counts as on because the header is present. If the header is
repeated, one value that is not `0` or `false` is enough.

A flagged request is served by a local model or a `trust: local` peer. If the
only target is remote, llama-swap does not call that peer. The response is
HTTP 503:

```json
{"src":"llama-swap","error":{"message":"rejected: local_only, no local target available","type":"server_error","param":null,"code":"service_unavailable"}}
```

The proxy log records `rejected: local_only, no local target available`. A
request without the header, or with `0` or `false`, is routed as before.

A model that no peer and no local config serves is still HTTP 404. The 503 is
only the refusal to use a remote target.

The header is forwarded on a peer hop. `Origin`, `Referer`, and `Sec-Fetch-*`
are stripped before the hop; `X-Local-Only` is not, so the next llama-swap
applies the same rule to its own peers. A `trust: local` peer that would
itself have to use a remote peer returns 503 rather than continuing onward.

A spillover selector can still name a remote target. The gate then rejects
that choice with 503. It does not skip back to an earlier local target.
