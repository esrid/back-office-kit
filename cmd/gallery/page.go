package main

import (
	"context"
	"io"
	"os"
)

const head = `<meta charset="utf-8">
<title>Back-Office Component Kit</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Bricolage+Grotesque:opsz,wght@12..96,600;12..96,700&family=IBM+Plex+Mono:wght@400;500&family=IBM+Plex+Sans:wght@400;500;600&display=swap">
`

// galleryCSS styles the documentation chrome only. Every colour is taken from
// a daisyUI token, so the page follows the viewer's theme in all three states
// (explicit light, explicit dark, and unstamped system).
const galleryCSS = `
:root { color-scheme: light dark; }
body {
  background: var(--color-base-100);
  color: var(--color-base-content);
  font-family: "IBM Plex Sans", ui-sans-serif, system-ui, sans-serif;
}
.gal {
  --g-rule: color-mix(in oklch, var(--color-base-content) 14%, transparent);
  --g-mute: color-mix(in oklch, var(--color-base-content) 62%, transparent);
  --g-faint: color-mix(in oklch, var(--color-base-content) 4%, transparent);
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  max-width: 1180px;
  margin-inline: auto;
  padding-inline: 1.25rem;
}
@media (min-width: 1000px) {
  .gal { grid-template-columns: 15rem minmax(0, 1fr); gap: 3.5rem; }
}

/* --- rail --- */
.gal-rail { display: none; }
@media (min-width: 1000px) {
  .gal-rail {
    display: block;
    position: sticky;
    top: 0;
    align-self: start;
    height: 100vh;
    padding-block: 3.5rem 2rem;
    overflow-y: auto;
  }
}
.gal-brand {
  display: flex; align-items: center; gap: .55rem;
  font-weight: 600; letter-spacing: -0.01em;
}
.gal-brand-mark {
  width: .7rem; height: .7rem; border-radius: 2px;
  background: var(--color-primary);
}
.gal-brand-sub {
  margin: .3rem 0 1rem 1.25rem;
  font-family: "IBM Plex Mono", ui-monospace, monospace;
  font-size: .7rem; color: var(--g-mute);
}
.gal-search {
  position: relative;
  margin-bottom: .6rem;
}
.gal-search input {
  width: 100%;
  padding: .38rem 2rem .38rem .6rem;
  font: inherit; font-size: .8rem;
  color: var(--color-base-content);
  background: color-mix(in oklch, var(--color-base-content) 4%, transparent);
  border: 1px solid var(--g-rule);
  border-radius: 6px;
  appearance: none;
}
.gal-search input:focus-visible {
  outline: 2px solid var(--color-primary);
  outline-offset: 1px;
  background: var(--color-base-100);
}
.gal-search input::-webkit-search-cancel-button { appearance: none; }
.gal-search kbd {
  position: absolute; right: .45rem; top: 50%; transform: translateY(-50%);
  font-family: "IBM Plex Mono", ui-monospace, monospace;
  font-size: .62rem; line-height: 1;
  padding: .18rem .3rem; border-radius: 3px;
  color: var(--g-mute);
  border: 1px solid var(--g-rule);
  pointer-events: none;
}
.gal-search input:focus-visible + kbd { display: none; }
.gal-count {
  font-family: "IBM Plex Mono", ui-monospace, monospace;
  font-size: .66rem; color: var(--g-mute);
  margin-bottom: 1rem;
}
.gal-empty {
  padding: 3rem 0;
  color: var(--g-mute);
  font-size: .9rem;
}

.gal-index { display: flex; flex-direction: column; }
.gal-index a {
  display: flex; gap: .7rem; align-items: baseline;
  padding: .34rem 0;
  font-size: .82rem; color: var(--g-mute);
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: color .12s ease;
}
.gal-index a:hover, .gal-index a:focus-visible { color: var(--color-base-content); }
.gal-index-n {
  font-family: "IBM Plex Mono", ui-monospace, monospace;
  font-size: .68rem;
  color: color-mix(in oklch, var(--color-base-content) 55%, transparent);
  font-variant-numeric: tabular-nums;
}
.gal-rail-foot {
  margin-top: 1.75rem; padding-top: 1rem;
  border-top: 1px solid var(--g-rule);
  font-size: .72rem; line-height: 1.5; color: var(--g-mute);
}

/* --- head --- */
.gal-main { padding-block: 3.5rem 5rem; min-width: 0; }
.gal-head { max-width: 44rem; margin-bottom: 4rem; }
.gal-eyebrow {
  font-family: "IBM Plex Mono", ui-monospace, monospace;
  font-size: .7rem; letter-spacing: .09em; text-transform: uppercase;
  color: var(--color-primary);
  margin-bottom: .9rem;
}
.gal-head h1 {
  font-family: "Bricolage Grotesque", ui-sans-serif, system-ui, sans-serif;
  font-weight: 700; font-size: clamp(1.9rem, 4vw, 2.75rem);
  line-height: 1.05; letter-spacing: -0.03em; text-wrap: balance;
  margin-bottom: .9rem;
}
.gal-lede { font-size: 1rem; line-height: 1.65; color: var(--g-mute); max-width: 38rem; }

/* --- section --- */
.gal-sec { padding-top: 2.75rem; margin-bottom: 3.5rem; border-top: 1px solid var(--g-rule); scroll-margin-top: 1.5rem; }
.gal-sec-head { display: flex; gap: 1rem; margin-bottom: 1.4rem; }
.gal-sec-meta { min-width: 0; }
.gal-sec-n {
  font-family: "IBM Plex Mono", ui-monospace, monospace;
  font-size: .72rem; padding-top: .35rem;
  color: color-mix(in oklch, var(--color-base-content) 55%, transparent);
  font-variant-numeric: tabular-nums;
}
.gal-sec h2 {
  font-family: "Bricolage Grotesque", ui-sans-serif, system-ui, sans-serif;
  font-weight: 600; font-size: 1.3rem; letter-spacing: -0.02em;
}
.gal-purpose { margin-top: .3rem; font-size: .93rem; line-height: 1.55; color: var(--g-mute); max-width: 40rem; }
.gal-uses { display: flex; flex-wrap: wrap; gap: .35rem; margin-top: .7rem; }
.gal-uses code {
  font-family: "IBM Plex Mono", ui-monospace, monospace;
  font-size: .68rem; padding: .12rem .4rem; border-radius: 3px;
  background: var(--g-faint); border: 1px solid var(--g-rule);
  color: var(--g-mute);
}

/* --- stage --- */
.gal-stage {
  border: 1px solid var(--g-rule);
  border-radius: 10px;
  background: var(--color-base-100);
  padding: 1.5rem;
  overflow-x: auto;
}
.gal-stage--app { padding: 0; overflow: hidden; }
.gal-stage--app .min-h-screen { min-height: 0; }
.gal-stage--app .drawer-content { min-height: 22rem; }
.gal-stage--app .drawer-side { height: 22rem; }

.gal-sig {
  margin-top: .85rem;
  font-family: "IBM Plex Mono", ui-monospace, monospace;
  font-size: .72rem; line-height: 1.6;
  color: var(--g-mute);
  overflow-x: auto;
  white-space: pre;
  padding-bottom: .2rem;
}
.gal-sig code { color: inherit; }

/* --- code + note --- */
.gal-code {
  margin-top: .85rem;
  padding: .95rem 1.1rem;
  border-radius: 8px;
  background: color-mix(in oklch, var(--color-base-content) 5%, transparent);
  border: 1px solid var(--g-rule);
  overflow-x: auto;
}
.gal-code code {
  font-family: "IBM Plex Mono", ui-monospace, monospace;
  font-size: .78rem; line-height: 1.65;
  white-space: pre;
  color: var(--color-base-content);
}
.gal-note {
  margin-top: .85rem; padding-left: .9rem;
  border-left: 2px solid color-mix(in oklch, var(--color-primary) 45%, transparent);
  font-size: .84rem; line-height: 1.6; color: var(--g-mute);
  max-width: 44rem;
}
.gal-foot {
  padding-top: 2rem; border-top: 1px solid var(--g-rule);
  font-size: .82rem; color: var(--g-mute);
}
.gal-foot code { font-family: "IBM Plex Mono", ui-monospace, monospace; }

@media (prefers-reduced-motion: reduce) {
  * { animation-duration: .01ms !important; transition-duration: .01ms !important; }
}
`

// writePage streams head, the compiled daisyUI stylesheet and the gallery body
// into w. The <title> is written first so publishers that scan only the head
// find it before the 75 KB stylesheet.
func writePage(w io.Writer, appCSSPath string, body func(io.Writer) error) error {
	css, err := os.ReadFile(appCSSPath)
	if err != nil {
		return err
	}
	if _, err := io.WriteString(w, head+"<style>"); err != nil {
		return err
	}
	if _, err := w.Write(css); err != nil {
		return err
	}
	if _, err := io.WriteString(w, galleryCSS+"</style>\n"); err != nil {
		return err
	}
	return body(w)
}

func renderGallery(w io.Writer) error {
	return Gallery(sections()).Render(context.Background(), w)
}
