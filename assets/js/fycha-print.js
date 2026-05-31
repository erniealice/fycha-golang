// fycha-print.js — delegated handler for print buttons.
//
// CSP-prep (Plan-6): the ledger general-ledger / trial-balance reports and the
// asset detail lapsing-schedule worksheet each carry a "Print" button that used
// to invoke `onclick="window.print()"`. Inline event-handler attributes cannot
// be nonced/hashed and block an enforcing `script-src 'self'`, so they are
// converted to document-level delegation via the shared `lf.on()` helper (see
// docs/wiki/articles/htmx-ui-patterns.md "JS event binding under hx-boost").
// Behavior is preserved exactly: clicking a `[data-print-trigger]` button opens
// the browser print dialog, identical to the prior inline `window.print()`.
//
// This is an APP-LEVEL fycha file (loaded per-page via `data-page-js`), NOT a
// pyeza-copied component. `lf-on.js` is loaded globally in app-shell.html before
// any page-js, so `lf.on` is always available when this runs. The page-js loader
// (app-shell.js) dedupes by src URL, so the handler registers exactly once per
// session even across HTMX navigations between the print pages.
(function () {
  'use strict';

  if (!window.lf || typeof window.lf.on !== 'function') {
    // lf-on.js must load before this file. If the helper is missing there is
    // nothing safe to bind to — bail quietly.
    return;
  }

  // Document-level delegation: matches `[data-print-trigger]` buttons even when
  // they are inside HTMX-swapped report/detail content.
  lf.on('click', '[data-print-trigger]', function () {
    window.print();
  });
})();
