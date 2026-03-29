/**
 * Journal Entry Form — client-side logic
 * Loaded in journal-list.html for the new journal entry sheet form.
 *
 * Handles:
 *   - Account search autocomplete via GET /action/ledger/accounts/search?q=...
 *   - Debit/credit mutual exclusion on each line
 *   - Running totals and balance validation
 *   - Post button enable/disable when balanced
 *   - Add line (appends a new row with next index)
 *   - Remove line (removes row; minimum 2 lines enforced)
 */

(function () {
  'use strict';

  var SEARCH_URL = '/action/ledger/accounts/search';
  var DEBOUNCE_MS = 300;
  var MIN_CHARS = 1;

  // ---------------------------------------------------------------------------
  // Autocomplete helpers
  // ---------------------------------------------------------------------------

  function attachAutocomplete(input) {
    if (input.dataset.acInit) return;
    input.dataset.acInit = '1';

    var lineIdx = input.dataset.line;
    var form = input.closest('form');
    var hiddenInput = form
      ? form.querySelector('input[name="account_id[' + lineIdx + ']"]')
      : null;

    var dropdown = document.createElement('ul');
    dropdown.className = 'je-ac-dropdown';
    dropdown.setAttribute('role', 'listbox');
    dropdown.style.cssText = [
      'position:absolute',
      'z-index:9999',
      'background:var(--color-surface,#fff)',
      'border:1px solid var(--color-border,#d1d5db)',
      'border-radius:var(--radius-md,6px)',
      'box-shadow:var(--shadow-md,0 4px 12px rgba(0,0,0,.15))',
      'margin-top:2px',
      'padding:4px 0',
      'list-style:none',
      'max-height:220px',
      'overflow-y:auto',
      'min-width:220px',
      'display:none',
    ].join(';');

    var wrapper = input.parentElement;
    if (wrapper && getComputedStyle(wrapper).position === 'static') {
      wrapper.style.position = 'relative';
    }
    input.after(dropdown);

    var debounceTimer = null;
    var abortController = null;
    var focusedIndex = -1;
    var lastResults = [];

    function openDropdown() { dropdown.style.display = 'block'; }
    function closeDropdown() { dropdown.style.display = 'none'; focusedIndex = -1; }

    function escapeHTML(str) {
      var d = document.createElement('div');
      d.textContent = str;
      return d.innerHTML;
    }

    function renderResults(results) {
      lastResults = results;
      focusedIndex = -1;
      dropdown.innerHTML = '';

      if (!results || results.length === 0) {
        var empty = document.createElement('li');
        empty.className = 'je-ac-empty';
        empty.style.cssText = 'padding:8px 12px;color:var(--color-text-muted,#6b7280);font-size:.875rem;';
        empty.textContent = 'No accounts found';
        dropdown.appendChild(empty);
        openDropdown();
        return;
      }

      results.forEach(function (item, i) {
        var li = document.createElement('li');
        li.className = 'je-ac-option';
        li.setAttribute('role', 'option');
        li.dataset.value = item.value;
        li.dataset.label = item.label;
        li.style.cssText = 'padding:7px 12px;cursor:pointer;font-size:.875rem;line-height:1.4;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;';
        li.innerHTML = escapeHTML(item.label);

        li.addEventListener('mouseenter', function () { setFocus(i); });
        li.addEventListener('mousedown', function (e) {
          e.preventDefault();
          selectOption(item);
        });
        dropdown.appendChild(li);
      });

      openDropdown();
    }

    function setFocus(idx) {
      var options = dropdown.querySelectorAll('.je-ac-option');
      options.forEach(function (el, i) {
        var focused = i === idx;
        el.style.background = focused ? 'var(--color-primary-50,#eff6ff)' : '';
        el.style.color = focused ? 'var(--color-primary-700,#1d4ed8)' : '';
      });
      focusedIndex = idx;
    }

    function selectOption(item) {
      input.value = item.label;
      if (hiddenInput) hiddenInput.value = item.value;
      closeDropdown();
      updateTotals();
    }

    function doSearch(term) {
      if (abortController) abortController.abort();
      abortController = new AbortController();
      var url = SEARCH_URL + '?q=' + encodeURIComponent(term);
      fetch(url, { signal: abortController.signal })
        .then(function (r) { return r.json(); })
        .then(renderResults)
        .catch(function (err) {
          if (err.name !== 'AbortError') closeDropdown();
        });
    }

    input.addEventListener('input', function () {
      var term = input.value.trim();
      if (hiddenInput) hiddenInput.value = '';
      clearTimeout(debounceTimer);
      if (term.length < MIN_CHARS) { closeDropdown(); return; }
      debounceTimer = setTimeout(function () { doSearch(term); }, DEBOUNCE_MS);
    });

    input.addEventListener('keydown', function (e) {
      if (dropdown.style.display === 'none') return;
      var options = dropdown.querySelectorAll('.je-ac-option');
      if (e.key === 'ArrowDown') { e.preventDefault(); setFocus(Math.min(focusedIndex + 1, options.length - 1)); }
      else if (e.key === 'ArrowUp') { e.preventDefault(); setFocus(Math.max(focusedIndex - 1, 0)); }
      else if (e.key === 'Enter') { e.preventDefault(); if (focusedIndex >= 0 && focusedIndex < lastResults.length) selectOption(lastResults[focusedIndex]); }
      else if (e.key === 'Escape') closeDropdown();
    });

    input.addEventListener('blur', function () { setTimeout(closeDropdown, 150); });
    input.addEventListener('focus', function () {
      var term = input.value.trim();
      if (term.length >= MIN_CHARS) doSearch(term);
    });
  }

  // ---------------------------------------------------------------------------
  // Balance / totals
  // ---------------------------------------------------------------------------

  function parseAmount(val) {
    if (!val || val === '') return 0;
    var n = parseFloat(String(val).replace(/,/g, ''));
    return isNaN(n) ? 0 : n;
  }

  function updateTotals() {
    var form = document.getElementById('journal-entry-form');
    if (!form) return;

    var totalDebit = 0;
    var totalCredit = 0;
    form.querySelectorAll('.je-debit').forEach(function (el) { totalDebit += parseAmount(el.value); });
    form.querySelectorAll('.je-credit').forEach(function (el) { totalCredit += parseAmount(el.value); });

    var diff = Math.abs(totalDebit - totalCredit);
    var balanced = diff < 0.005 && totalDebit > 0;

    var alertEl = document.getElementById('je-balance-alert');
    var msgEl = document.getElementById('je-balance-message');
    var postBtn = document.getElementById('je-post-btn');

    if (balanced) {
      if (alertEl) alertEl.className = 'alert alert--success je-balance-alert';
      if (msgEl) msgEl.textContent = 'Balanced \u2014 Debits and credits are equal (' + totalDebit.toFixed(2) + ')';
      if (postBtn) postBtn.disabled = false;
    } else {
      if (alertEl) alertEl.className = (totalDebit === 0 && totalCredit === 0) ? 'alert alert--info je-balance-alert' : 'alert alert--danger je-balance-alert';
      if (msgEl) {
        if (totalDebit === 0 && totalCredit === 0) {
          msgEl.textContent = 'Enter debit and credit amounts to balance the journal entry.';
        } else {
          msgEl.textContent = 'Unbalanced \u2014 Debits: ' + totalDebit.toFixed(2) + ' / Credits: ' + totalCredit.toFixed(2) + ' (difference: ' + diff.toFixed(2) + ')';
        }
      }
      if (postBtn) postBtn.disabled = true;
    }
  }

  // ---------------------------------------------------------------------------
  // Debit / Credit mutual exclusion
  // ---------------------------------------------------------------------------

  function attachAmountHandlers(row) {
    var debit = row.querySelector('.je-debit');
    var credit = row.querySelector('.je-credit');
    if (debit) debit.addEventListener('input', function () { if (parseAmount(debit.value) > 0 && credit) credit.value = ''; updateTotals(); });
    if (credit) credit.addEventListener('input', function () { if (parseAmount(credit.value) > 0 && debit) debit.value = ''; updateTotals(); });
  }

  // ---------------------------------------------------------------------------
  // Add / Remove line
  // ---------------------------------------------------------------------------

  function getNextLineIndex(form) {
    var max = 0;
    form.querySelectorAll('.journal-line').forEach(function (r) {
      var idx = parseInt(r.dataset.line, 10);
      if (idx > max) max = idx;
    });
    return max + 1;
  }

  function buildNewRow(idx) {
    var tr = document.createElement('tr');
    tr.className = 'journal-line';
    tr.dataset.line = String(idx);
    tr.innerHTML =
      '<td class="col-num">' + idx + '</td>' +
      '<td class="col-account">' +
        '<input type="hidden" name="account_id[' + idx + ']" value="" class="je-account-id">' +
        '<input type="text" id="je-account-' + idx + '" class="form-control je-account-search"' +
          ' placeholder="Search account\u2026" aria-label="Account, line ' + idx + '"' +
          ' autocomplete="off" data-line="' + idx + '">' +
      '</td>' +
      '<td class="col-debit">' +
        '<input type="number" id="je-debit-' + idx + '" name="debit[' + idx + ']"' +
          ' class="form-control je-debit" value="" min="0" step="0.01"' +
          ' placeholder="0.00" aria-label="Debit, line ' + idx + '" data-line="' + idx + '">' +
      '</td>' +
      '<td class="col-credit">' +
        '<input type="number" id="je-credit-' + idx + '" name="credit[' + idx + ']"' +
          ' class="form-control je-credit" value="" min="0" step="0.01"' +
          ' placeholder="0.00" aria-label="Credit, line ' + idx + '" data-line="' + idx + '">' +
      '</td>' +
      '<td class="col-memo">' +
        '<input type="text" id="je-memo-' + idx + '" name="memo[' + idx + ']"' +
          ' class="form-control je-memo" value="" aria-label="Memo, line ' + idx + '" placeholder="">' +
      '</td>' +
      '<td class="col-remove">' +
        '<button type="button" class="btn-icon btn-icon--danger je-remove-line"' +
          ' data-line="' + idx + '" title="Remove line">' +
          '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"' +
            ' stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">' +
            '<line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>' +
          '</svg>' +
        '</button>' +
      '</td>';
    return tr;
  }

  function wireRow(row) {
    var searchInput = row.querySelector('.je-account-search');
    if (searchInput) attachAutocomplete(searchInput);
    attachAmountHandlers(row);
    var removeBtn = row.querySelector('.je-remove-line');
    if (removeBtn) {
      removeBtn.addEventListener('click', function () {
        var form = document.getElementById('journal-entry-form');
        if (!form) return;
        if (form.querySelectorAll('.journal-line').length <= 2) return;
        row.remove();
        updateTotals();
      });
    }
  }

  // ---------------------------------------------------------------------------
  // Init
  // ---------------------------------------------------------------------------

  function initJournalForm() {
    var form = document.getElementById('journal-entry-form');
    if (!form || form.dataset.jeInit) return;
    form.dataset.jeInit = '1';

    form.querySelectorAll('.journal-line').forEach(wireRow);
    updateTotals();

    var addBtn = document.getElementById('je-add-line');
    if (addBtn) {
      addBtn.addEventListener('click', function () {
        var tbody = document.getElementById('journal-lines-body');
        if (!tbody) return;
        var idx = getNextLineIndex(form);
        var newRow = buildNewRow(idx);
        tbody.appendChild(newRow);
        wireRow(newRow);
        var newInput = newRow.querySelector('.je-account-search');
        if (newInput) setTimeout(function () { newInput.focus(); }, 50);
      });
    }
  }

  function tryInit() {
    if (document.getElementById('journal-entry-form')) initJournalForm();
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', tryInit);
  } else {
    tryInit();
  }

  document.addEventListener('htmx:afterSwap', tryInit);

})();
