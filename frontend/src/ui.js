// This module stores the UI functionality for the frontend

export function initUI() {
    const app = document.getElementById('app');
    app.innerHTML = `
    <!-- Noise canvas background -->
    <canvas id="noise-canvas"></canvas>

    <!-- Grid overlay -->
    <div class="grid-overlay"></div>

    <!-- Header -->
    <header class="site-header">
      <div class="logo">
        <span class="logo-bracket">[</span>
        <span class="logo-text">Phake<span class="logo-accent">Links</span></span>
        <span class="logo-bracket">]</span>
      </div>
      <p class="tagline">// dissect deception. stay protected.</p>
    </header>

    <!-- Main content -->
    <main class="main-content">

      <!-- Mode Toggle -->
      <div class="mode-toggle-wrapper">
        <span class="mode-label">MODE://</span>
        <div class="mode-toggle" id="modeToggle" role="button" tabindex="0" aria-label="Toggle mode">
          <div class="toggle-option active" data-mode="educational" id="modeEdu">
            <span class="toggle-icon">⬡</span> EDUCATIONAL
          </div>
          <div class="toggle-option" data-mode="prank" id="modePrank">
            <span class="toggle-icon">⚠</span> PRANK
          </div>
          <div class="toggle-slider" id="toggleSlider"></div>
        </div>
      </div>

      <!-- Description blurb that changes with mode -->
      <p class="mode-description" id="modeDesc">
        Analyze any URL to learn how phishing links are crafted — see the technique, understand the trick.
      </p>

      <!-- Input area -->
      <div class="input-section">
        <div class="input-wrapper" id="inputWrapper">
          <span class="input-prefix">URL://</span>
          <input
            type="text"
            id="urlInput"
            class="url-input"
            placeholder="paste link, domain or url..."
            autocomplete="off"
            spellcheck="false"
          />
          <div class="input-scan-line"></div>
        </div>

        <button class="submit-btn" id="submitBtn" type="button">
          <span class="btn-text">ANALYZE</span>
          <span class="btn-arrow">→</span>
          <div class="btn-glitch"></div>
        </button>
      </div>

      <!-- Corner decorations -->
      <div class="corner-decor top-left"></div>
      <div class="corner-decor top-right"></div>
      <div class="corner-decor bottom-left"></div>
      <div class="corner-decor bottom-right"></div>

    </main>

    <!-- Loading Overlay -->
    <div class="overlay" id="loadingOverlay">
      <div class="loading-box">
        <div class="loading-spinner">
          <div class="spinner-ring"></div>
          <div class="spinner-ring ring2"></div>
          <div class="spinner-ring ring3"></div>
        </div>
        <p class="loading-text" id="loadingText">SCANNING TARGET...</p>
        <div class="loading-bar-track">
          <div class="loading-bar-fill"></div>
        </div>
        <p class="loading-sub">analyzing link structure</p>
      </div>
    </div>

    <!-- Educational Result Modal -->
    <div class="overlay" id="eduModal">
      <div class="modal edu-modal">
        <button class="modal-close" id="closeEduModal" aria-label="Close">✕</button>

        <div class="modal-header">
          <span class="modal-tag">// ANALYSIS COMPLETE</span>
          <h2 class="modal-title">Phishing Breakdown</h2>
        </div>

        <div class="edu-grid">
          <div class="edu-card" id="originalLinkCard">
            <span class="card-label">ORIGINAL LINK</span>
            <p class="card-value original-link" id="originalLink">—</p>
          </div>

          <div class="edu-card highlight" id="fakeLinkCard">
            <span class="card-label">GENERATED FAKE</span>
            <p class="card-value fake-link" id="fakeLink">—</p>
            <button class="copy-btn small" data-target="fakeLink">COPY</button>
          </div>

          <div class="edu-card full-width" id="techniqueCard">
            <span class="card-label">TECHNIQUE USED</span>
            <p class="card-value technique-tag" id="techniqueUsed">—</p>
          </div>

          <div class="edu-card full-width explanation-card" id="explanationCard">
            <span class="card-label">EXPLANATION</span>
            <p class="card-value explanation-text" id="explanationText">—</p>
          </div>
        </div>

        <div class="modal-footer">
          <span class="footer-note">// knowledge is your firewall</span>
        </div>
      </div>
    </div>

    <!-- Prank Result Modal -->
    <div class="overlay" id="prankModal">
      <div class="modal prank-modal">
        <button class="modal-close prank-close" id="closePrankModal" aria-label="Close">✕</button>

        <div class="modal-header">
          <span class="modal-tag prank-tag">// PRANK LINK GENERATED</span>
          <h2 class="modal-title prank-title">Send it. 😈</h2>
        </div>

        <div class="prank-link-display">
          <div class="prank-warning-strip">⚠ LOOKS SKETCHY ⚠ PROCEED WITH CAUTION ⚠ LOOKS SKETCHY ⚠ PROCEED WITH CAUTION ⚠</div>
          <div class="prank-link-box">
            <p class="prank-link-text" id="prankLinkText">—</p>
          </div>
          <button class="copy-btn prank-copy-btn" id="prankCopyBtn">
            <span id="prankCopyLabel">COPY LINK</span>
          </button>
        </div>

        <p class="prank-disclaimer">For entertainment only. Don't be evil.</p>
      </div>
    </div>

    <!-- Error Modal -->
    <div class="overlay" id="errorModal">
      <div class="modal error-modal">
        <button class="modal-close" id="closeErrorModal" aria-label="Close">✕</button>
        <div class="error-icon">!</div>
        <span class="modal-tag error-tag">// ERROR DETECTED</span>
        <p class="error-message" id="errorMessage">Something went wrong.</p>
        <button class="submit-btn error-retry-btn" id="retryBtn">RETRY</button>
      </div>
    </div>
  `;

    // Init noise background
    initNoise();

    // Init interactions
    setupModeToggle();
    setupInputEffects();
    setupModalClosers();
    setupCopyButtons();
}

// Mode Toggle

let currentMode = 'educational';

function setupModeToggle() {
    const toggle = document.getElementById('modeToggle');
    const eduOpt = document.getElementById('modeEdu');
    const prankOpt = document.getElementById('modePrank');
    const desc = document.getElementById('modeDesc');
    const submitBtn = document.getElementById('submitBtn');
    const inputWrapper = document.getElementById('inputWrapper');

    const descriptions = {
        educational:
            'Analyze any URL to learn how phishing links are crafted — see the technique, understand the trick.',
        prank: 'Generate a hilariously suspicious fake link to send to your friends. Harmless chaos.',
    };

    const btnLabels = { educational: 'ANALYZE', prank: 'GENERATE' };

    function setMode(mode) {
        currentMode = mode;
        document.body.dataset.mode = mode;

        eduOpt.classList.toggle('active', mode === 'educational');
        prankOpt.classList.toggle('active', mode === 'prank');

        desc.classList.add('fade-out');
        setTimeout(() => {
            desc.textContent = descriptions[mode];
            desc.classList.remove('fade-out');
        }, 200);

        submitBtn.querySelector('.btn-text').textContent = btnLabels[mode];
        inputWrapper.dataset.mode = mode;
    }

    eduOpt.addEventListener('click', () => setMode('educational'));
    prankOpt.addEventListener('click', () => setMode('prank'));
    toggle.addEventListener('keydown', (e) => {
        if (e.key === 'ArrowRight') setMode('prank');
        if (e.key === 'ArrowLeft') setMode('educational');
    });

    document.body.dataset.mode = 'educational';
}

export function getMode() {
    return currentMode;
}

export function getlink() {
    return document.getElementById('urlInput')?.value?.trim() || '';
}

// Input Effects

function setupInputEffects() {
    const input = document.getElementById('urlInput');
    const wrapper = document.getElementById('inputWrapper');

    input.addEventListener('focus', () => wrapper.classList.add('focused'));
    input.addEventListener('blur', () => wrapper.classList.remove('focused'));
    input.addEventListener('input', () => {
        wrapper.classList.toggle('has-value', input.value.length > 0);
    });

    input.addEventListener('keydown', (e) => {
        if (e.key === 'Enter') {
            document.getElementById('submitBtn').click();
        }
    });

    // Button hover glitch
    const btn = document.getElementById('submitBtn');
    btn.addEventListener('mouseenter', () => btn.classList.add('glitch'));
    btn.addEventListener('mouseleave', () => btn.classList.remove('glitch'));
}

// Loading

const loadingMessages = {
    educational: [
        'SCANNING TARGET...',
        'PARSING STRUCTURE...',
        'DETECTING TECHNIQUE...',
        'BUILDING REPORT...',
    ],
    prank: [
        'CRAFTING CHAOS...',
        'MAXIMUM SKETCHINESS...',
        'ADDING RED FLAGS...',
        'ALMOST EVIL ENOUGH...',
    ],
};

let loadingInterval = null;

export function showLoading() {
    const overlay = document.getElementById('loadingOverlay');
    const loadingText = document.getElementById('loadingText');
    const messages =
        loadingMessages[currentMode] || loadingMessages.educational;
    let i = 0;

    overlay.classList.add('active');
    loadingText.textContent = messages[0];

    loadingInterval = setInterval(() => {
        i = (i + 1) % messages.length;
        loadingText.textContent = messages[i];
    }, 900);
}

export function hideLoading() {
    const overlay = document.getElementById('loadingOverlay');
    overlay.classList.remove('active');
    if (loadingInterval) clearInterval(loadingInterval);
}

// Edu Result

export function showEducationalResult(data) {
    hideLoading();
    // data = { originalLink, fakeLink, technique, explanation }
    document.getElementById('originalLink').textContent =
        data.originalLink || '—';
    document.getElementById('fakeLink').textContent = data.fakeLink || '—';
    document.getElementById('techniqueUsed').textContent =
        data.technique || '—';
    document.getElementById('explanationText').textContent =
        data.explanation || '—';

    // Staggered card reveal
    const cards = document.querySelectorAll('.edu-card');
    cards.forEach((c, i) => {
        c.style.opacity = 0;
        c.style.transform = 'translateY(16px)';
        setTimeout(
            () => {
                c.style.transition = 'opacity 0.4s ease, transform 0.4s ease';
                c.style.opacity = 1;
                c.style.transform = 'translateY(0)';
            },
            100 + i * 100,
        );
    });

    document.getElementById('eduModal').classList.add('active');
}

// Prank Result

export function showPrankResult(link) {
    hideLoading();
    document.getElementById('prankLinkText').textContent = link || '—';
    document.getElementById('prankModal').classList.add('active');
}

// Error

export function showError(msg) {
    hideLoading();
    document.getElementById('errorMessage').textContent =
        msg || 'An unexpected error occurred.';
    document.getElementById('errorModal').classList.add('active');
}

// Modal Closers

function setupModalClosers() {
    const pairs = [
        ['closeEduModal', 'eduModal'],
        ['closePrankModal', 'prankModal'],
        ['closeErrorModal', 'errorModal'],
        ['retryBtn', 'errorModal'],
    ];

    pairs.forEach(([btnId, modalId]) => {
        document.getElementById(btnId)?.addEventListener('click', () => {
            document.getElementById(modalId).classList.remove('active');
        });
    });

    // Close on overlay click
    ['eduModal', 'prankModal', 'errorModal'].forEach((id) => {
        const overlay = document.getElementById(id);
        overlay.addEventListener('click', (e) => {
            if (e.target === overlay) overlay.classList.remove('active');
        });
    });

    // ESC key
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            ['eduModal', 'prankModal', 'errorModal'].forEach((id) => {
                document.getElementById(id)?.classList.remove('active');
            });
        }
    });
}

// Copy Buttons

function setupCopyButtons() {
    document.addEventListener('click', (e) => {
        // Small copy buttons in edu modal
        if (e.target.closest('.copy-btn.small')) {
            const targetId = e.target.closest('.copy-btn').dataset.target;
            const text = document.getElementById(targetId)?.textContent;
            copyToClipboard(text, e.target.closest('.copy-btn'));
        }

        // Prank copy btn
        if (e.target.closest('#prankCopyBtn')) {
            const text = document.getElementById('prankLinkText')?.textContent;
            copyToClipboard(
                text,
                document.getElementById('prankCopyBtn'),
                'prankCopyLabel',
                'COPIED! 🫣',
            );
        }
    });
}

function copyToClipboard(text, btn, labelId = null, successMsg = 'COPIED!') {
    if (!text || text === '—') return;
    navigator.clipboard.writeText(text).then(() => {
        const original = labelId
            ? document.getElementById(labelId)?.textContent
            : btn.textContent;
        if (labelId) document.getElementById(labelId).textContent = successMsg;
        else btn.textContent = successMsg;
        btn.classList.add('copied');
        setTimeout(() => {
            if (labelId)
                document.getElementById(labelId).textContent = original;
            else btn.textContent = original;
            btn.classList.remove('copied');
        }, 2000);
    });
}

// Noise Canvas

function initNoise() {
    const canvas = document.getElementById('noise-canvas');
    const ctx = canvas.getContext('2d');
    let frame = 0;

    function resize() {
        canvas.width = window.innerWidth;
        canvas.height = window.innerHeight;
    }

    function drawNoise() {
        const w = canvas.width,
            h = canvas.height;
        const imageData = ctx.createImageData(w, h);
        const data = imageData.data;

        for (let i = 0; i < data.length; i += 4) {
            const val = Math.random() * 12;
            data[i] = val;
            data[i + 1] = val;
            data[i + 2] = val;
            data[i + 3] = 18;
        }

        ctx.putImageData(imageData, 0, 0);
        frame++;
        if (frame % 3 === 0) requestAnimationFrame(drawNoise);
        else setTimeout(() => requestAnimationFrame(drawNoise), 80);
    }

    resize();
    window.addEventListener('resize', resize);
    drawNoise();
}
