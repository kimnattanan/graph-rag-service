import * as api from './api.js';

const sessionKey = 'grs.session';

const state = {
  user: null,
  view: 'login',
  banner: '',
  bannerKind: '',
  conversations: null,
  conversation: null,
  selectedConversationId: null,
  documents: null,
  document: null,
  indexJob: null,
  selectedDocumentId: null,
  docFilter: { tag: '', indexStatus: '' },
  docDraft: { title: '', tags: '', content: '' },
  draft: '',
  sending: false,
};

const viewEl = document.querySelector('#view');
const navEl = document.querySelector('#nav');
const whoEl = document.querySelector('#who');
const bannerEl = document.querySelector('#banner');

function isAdmin() {
  return state.user?.role === 'admin' || state.user?.permissions?.includes('knowledge:write');
}

function esc(value) {
  return String(value ?? '')
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;');
}

function roleLabel(role) {
  if (role === 'user') return 'You';
  if (role === 'assistant') return 'Assistant';
  if (!role) return 'Message';
  const text = String(role);
  return text.charAt(0).toUpperCase() + text.slice(1);
}

function statusLabel(status) {
  if (!status) return 'Any status';
  const text = String(status);
  return text.charAt(0).toUpperCase() + text.slice(1);
}

function formatWhen(value) {
  if (!value) return '—';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return String(value);
  return new Intl.DateTimeFormat(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date);
}

function formatScore(score) {
  const number = Number(score);
  if (!Number.isFinite(number)) return String(score ?? '');
  return number.toFixed(3);
}

function navLink(hash, label) {
  const active = state.view === hash;
  return `<a class="nav-link${active ? ' is-active' : ''}" href="#${hash}"${active ? ' aria-current="page"' : ''}>${label}</a>`;
}

function parseTags(value) {
  return value.split(',').map((part) => part.trim()).filter(Boolean);
}

function currentHash() {
  return location.hash.replace(/^#/, '');
}

function setBanner(kind, message) {
  state.banner = message;
  state.bannerKind = kind;
}

function fail(err) {
  if (err.status === 401 && state.user) {
    clearSession();
    setBanner('error', err.slug ? `${err.slug}: ${err.message}` : err.message);
    if (currentHash() !== 'login') {
      location.hash = 'login';
      return;
    }
  } else {
    setBanner('error', err.slug ? `${err.slug}: ${err.message}` : err.message || String(err));
  }
  render();
}

function loadSession() {
  const raw = localStorage.getItem(sessionKey);
  if (!raw) return;
  try {
    const session = JSON.parse(raw);
    api.setToken(session.accessToken);
    state.user = session.user;
  } catch {
    localStorage.removeItem(sessionKey);
  }
}

function saveSession(auth) {
  api.setToken(auth.accessToken);
  state.user = auth.user;
  localStorage.setItem(sessionKey, JSON.stringify({
    accessToken: auth.accessToken,
    user: auth.user,
  }));
}

function clearSession() {
  api.setToken(null);
  state.user = null;
  state.conversations = null;
  state.conversation = null;
  state.selectedConversationId = null;
  state.documents = null;
  state.document = null;
  state.indexJob = null;
  state.selectedDocumentId = null;
  localStorage.removeItem(sessionKey);
}

function render() {
  const links = [];
  if (state.user) {
    links.push(navLink('chat', 'Conversation'));
    if (isAdmin()) links.push(navLink('documents', 'Documents'));
  } else {
    links.push(navLink('login', 'Sign in'));
    links.push(navLink('register', 'Register'));
  }
  navEl.innerHTML = links.join('');

  if (state.user) {
    whoEl.innerHTML = `
      <span class="account-name">${esc(state.user.username)}</span>
      <span class="role-badge">${esc(state.user.role)}</span>
      <button type="button" class="btn btn-ghost" id="logout">Sign out</button>
      <button type="button" class="btn btn-danger" id="delete-account">Delete account</button>
    `;
  } else {
    whoEl.innerHTML = '';
  }

  if (state.banner) {
    bannerEl.hidden = false;
    bannerEl.className = `banner ${state.bannerKind}`;
    bannerEl.textContent = state.banner;
  } else {
    bannerEl.hidden = true;
    bannerEl.className = 'banner';
    bannerEl.textContent = '';
  }

  if (state.view === 'login') viewEl.innerHTML = loginView();
  else if (state.view === 'register') viewEl.innerHTML = registerView();
  else if (state.view === 'chat') viewEl.innerHTML = chatView();
  else if (state.view === 'documents') viewEl.innerHTML = documentsView();
  else viewEl.innerHTML = '<p>Unknown page.</p>';

  bind();
  const messages = viewEl.querySelector('.messages');
  if (messages) messages.scrollTop = messages.scrollHeight;
}

function loginView() {
  return `
    <section class="sheet">
      <h1>Sign in</h1>
      <p class="lede">Use your account to open a conversation. Administrators can also manage documents.</p>
      <form id="login-form">
        <label class="field">
          <span>Email</span>
          <input name="email" type="email" required autocomplete="username">
        </label>
        <label class="field">
          <span>Password</span>
          <input name="password" type="password" required autocomplete="current-password">
        </label>
        <button class="btn">Sign in</button>
      </form>
      <p class="note">Local services: user on port 3001, knowledge on port 3000, conversation on port 3002.</p>
    </section>
  `;
}

function registerView() {
  return `
    <section class="sheet">
      <h1>Register</h1>
      <p class="lede">Create an account, then sign in to continue.</p>
      <form id="register-form">
        <label class="field">
          <span>Email</span>
          <input name="email" type="email" required autocomplete="email">
        </label>
        <label class="field">
          <span>Username</span>
          <input name="username" type="text" required maxlength="64" autocomplete="username">
        </label>
        <label class="field">
          <span>Password</span>
          <input name="password" type="password" required minlength="8" maxlength="72" autocomplete="new-password">
          <small class="hint">At least 8 characters.</small>
        </label>
        <button class="btn">Create account</button>
      </form>
    </section>
  `;
}

function chatView() {
  const items = state.conversations?.items || [];
  const list = state.conversations
    ? items.map((item) => {
        const selected = item.id === state.selectedConversationId;
        const count = Number(item.messageCount) || 0;
        return `
          <div class="conv-item${selected ? ' is-selected' : ''}">
            <button type="button" class="conv-open" data-open-conversation="${esc(item.id)}" ${selected ? 'disabled aria-current="true"' : ''}>
              <span class="conv-title">${esc(item.title || 'Untitled')}</span>
              <span class="conv-meta">${count} ${count === 1 ? 'message' : 'messages'}</span>
            </button>
            <button type="button" class="btn-text" data-delete-conversation="${esc(item.id)}">Delete</button>
          </div>
        `;
      }).join('') || '<p class="empty">No conversations yet.</p>'
    : `<p class="empty">${state.banner ? 'Could not load conversations.' : 'Loading…'}</p>`;

  const messages = state.conversation?.messages?.map((message) => {
    const role = message.role === 'assistant' || message.role === 'user' ? message.role : 'other';
    const sources = (message.sources || []).map((source) => `
      <li>
        <div class="source-title">${esc(source.documentTitle)}</div>
        <div class="source-score">Score ${esc(formatScore(source.score))}</div>
        <div>${esc(source.text)}</div>
      </li>
    `).join('');
    return `
      <article class="message message-${esc(role)}">
        <div class="role">${esc(roleLabel(message.role))}</div>
        <div class="message-body">${esc(message.content)}</div>
        ${sources ? `<details class="sources"><summary>Sources (${message.sources.length})</summary><ul class="source-list">${sources}</ul></details>` : ''}
      </article>
    `;
  }).join('') || (state.conversation
    ? '<p class="empty">No messages yet. Write a question below.</p>'
    : '<p class="empty">Select a conversation, or create one to begin.</p>');

  const total = state.conversations ? state.conversations.total : null;

  return `
    <div class="chat">
      <aside class="side">
        <div class="side-head">
          <h2>Conversations</h2>
          <form id="new-conversation">
            <label class="field">
              <span>Title</span>
              <input name="title" type="text" placeholder="Optional">
            </label>
            <button class="btn">New conversation</button>
          </form>
        </div>
        ${total === null ? '' : `<p class="count">${total} ${total === 1 ? 'conversation' : 'conversations'}</p>`}
        <div class="conv-list">${list}</div>
      </aside>
      <section class="thread">
        <header class="thread-head">
          <h2>${esc(state.conversation?.title || 'Conversation')}</h2>
        </header>
        <div class="messages">${messages}</div>
        <form id="composer" class="composer">
          <label class="field">
            <span>Message</span>
            <textarea name="content" required ${state.sending ? 'disabled' : ''}>${esc(state.draft)}</textarea>
          </label>
          <div class="composer-options">
            <label class="field">
              <span>Passages</span>
              <input name="topK" type="number" min="1" max="50" value="5">
            </label>
            <label class="field">
              <span>Prior messages</span>
              <input name="historyCapacity" type="number" min="0" max="100" value="3">
            </label>
            <label class="field">
              <span>Tags</span>
              <input name="tags" type="text" placeholder="Comma-separated">
            </label>
          </div>
          <div class="actions">
            <button class="btn" ${state.sending || !state.conversation ? 'disabled' : ''}>${state.sending ? 'Sending…' : 'Send message'}</button>
          </div>
        </form>
      </section>
    </div>
  `;
}

function documentsView() {
  if (!isAdmin()) {
    return `
      <section class="sheet">
        <h1>Documents</h1>
        <p class="lede">Document management requires the knowledge:write permission.</p>
      </section>
    `;
  }

  const items = state.documents?.items || [];
  const rows = items.map((item) => `
    <tr class="${item.id === state.selectedDocumentId ? 'selected' : ''}">
      <td>${esc(item.title)}</td>
      <td>${esc((item.tags || []).join(', ') || '—')}</td>
      <td><span class="status status-${esc(item.indexStatus)}">${esc(statusLabel(item.indexStatus))}</span></td>
      <td class="nowrap">${esc(formatWhen(item.updatedAt))}</td>
      <td class="nowrap">
        <div class="row-actions">
          <button type="button" class="btn-quiet" data-open-document="${esc(item.id)}">Open</button>
          <button type="button" class="btn-text" data-delete-document="${esc(item.id)}">Delete</button>
        </div>
      </td>
    </tr>
  `).join('');

  const doc = state.document;
  const job = state.indexJob;
  const editor = doc ? `
    <form id="edit-document">
      <h2>Edit document</h2>
      <p class="meta-line"><span class="mono">${esc(doc.id)}</span> · <span class="status status-${esc(doc.indexStatus)}">${esc(statusLabel(doc.indexStatus))}</span></p>
      <label class="field">
        <span>Title</span>
        <input name="title" type="text" required value="${esc(doc.title)}">
      </label>
      <label class="field">
        <span>Tags</span>
        <input name="tags" type="text" value="${esc((doc.tags || []).join(', '))}" placeholder="Comma-separated">
      </label>
      <label class="field">
        <span>Content</span>
        <textarea name="content" required>${esc(doc.content)}</textarea>
      </label>
      <div class="actions">
        <button class="btn">Save</button>
        <button type="button" class="btn btn-ghost" id="reindex">Reindex</button>
        <button type="button" class="btn btn-ghost" id="index-status">Index status</button>
        <button type="button" class="btn btn-danger" id="delete-document">Delete</button>
      </div>
    </form>
    ${job ? `
      <dl class="job">
        <dt>Status</dt><dd>${esc(statusLabel(job.status))}</dd>
        <dt>Started</dt><dd>${esc(formatWhen(job.startedAt))}</dd>
        <dt>Finished</dt><dd>${esc(formatWhen(job.finishedAt))}</dd>
        <dt>Error</dt><dd>${esc(job.errorMessage || '—')}</dd>
      </dl>
    ` : ''}
  ` : `
    <div class="editor-empty">
      <h2>Document</h2>
      <p>Select a document to review its content, save changes, or reindex it.</p>
    </div>
  `;

  const total = state.documents ? state.documents.total : null;

  return `
    <div class="docs-page">
      <header class="page-head">
        <h1>Documents</h1>
        <p class="lede">Source documents are stored, indexed, and retrieved as passages in the knowledge graph.</p>
      </header>
      <form id="doc-filter" class="toolbar">
        <label class="field">
          <span>Tag</span>
          <input name="tag" type="text" value="${esc(state.docFilter.tag)}">
        </label>
        <label class="field">
          <span>Index status</span>
          <select name="indexStatus">
            ${['', 'pending', 'indexing', 'completed', 'failed'].map((status) => `
              <option value="${status}" ${state.docFilter.indexStatus === status ? 'selected' : ''}>${esc(statusLabel(status))}</option>
            `).join('')}
          </select>
        </label>
        <button class="btn">Apply</button>
      </form>
      <div class="docs">
        <section class="card card-flush">
          <div class="card-head">
            <h2>Library</h2>
            ${total === null ? '' : `<p class="count">${total} ${total === 1 ? 'document' : 'documents'}</p>`}
          </div>
          <div class="doc-list">
            ${state.documents ? `
              <table>
                <thead>
                  <tr>
                    <th>Title</th>
                    <th>Tags</th>
                    <th>Index</th>
                    <th>Updated</th>
                    <th></th>
                  </tr>
                </thead>
                <tbody>${rows || '<tr><td colspan="5">No documents match this filter.</td></tr>'}</tbody>
              </table>
            ` : `<p class="empty">${state.banner ? 'Could not load documents.' : 'Loading…'}</p>`}
          </div>
        </section>
        <section class="card editor">${editor}</section>
      </div>
      <section class="card">
        <form id="create-document">
          <h2>New document</h2>
          <label class="field">
            <span>Title</span>
            <input name="title" type="text" required value="${esc(state.docDraft.title)}">
          </label>
          <label class="field">
            <span>Tags</span>
            <input name="tags" type="text" placeholder="Comma-separated" value="${esc(state.docDraft.tags)}">
          </label>
          <label class="field">
            <span>Content</span>
            <textarea name="content" required placeholder="Markdown">${esc(state.docDraft.content)}</textarea>
          </label>
          <button class="btn">Create document</button>
        </form>
      </section>
    </div>
  `;
}

function bind() {
  document.querySelector('#logout')?.addEventListener('click', onLogout);
  document.querySelector('#delete-account')?.addEventListener('click', onDeleteAccount);
  document.querySelector('#login-form')?.addEventListener('submit', onLogin);
  document.querySelector('#register-form')?.addEventListener('submit', onRegister);
  document.querySelector('#new-conversation')?.addEventListener('submit', onCreateConversation);
  document.querySelector('#composer')?.addEventListener('submit', onSend);
  document.querySelector('#doc-filter')?.addEventListener('submit', onFilterDocuments);
  document.querySelector('#create-document')?.addEventListener('submit', onCreateDocument);
  document.querySelector('#edit-document')?.addEventListener('submit', onUpdateDocument);
  document.querySelector('#reindex')?.addEventListener('click', onReindex);
  document.querySelector('#index-status')?.addEventListener('click', onIndexStatus);
  document.querySelector('#delete-document')?.addEventListener('click', () => {
    if (state.selectedDocumentId) void removeDocument(state.selectedDocumentId);
  });

  viewEl.querySelectorAll('[data-open-conversation]').forEach((button) => {
    button.addEventListener('click', () => openConversation(button.dataset.openConversation));
  });
  viewEl.querySelectorAll('[data-delete-conversation]').forEach((button) => {
    button.addEventListener('click', () => removeConversation(button.dataset.deleteConversation));
  });
  viewEl.querySelectorAll('[data-open-document]').forEach((button) => {
    button.addEventListener('click', () => openDocument(button.dataset.openDocument));
  });
  viewEl.querySelectorAll('[data-delete-document]').forEach((button) => {
    button.addEventListener('click', () => removeDocument(button.dataset.deleteDocument));
  });
}

async function onLogin(event) {
  event.preventDefault();
  const data = new FormData(event.target);
  try {
    const auth = await api.login({
      email: data.get('email'),
      password: data.get('password'),
    });
    saveSession(auth);
    setBanner('ok', `Logged in as ${auth.user.username}`);
    location.hash = 'chat';
  } catch (err) {
    fail(err);
  }
}

async function onRegister(event) {
  event.preventDefault();
  const data = new FormData(event.target);
  try {
    await api.register({
      email: data.get('email'),
      username: data.get('username'),
      password: data.get('password'),
    });
    setBanner('ok', 'Account created. Log in to continue.');
    location.hash = 'login';
  } catch (err) {
    fail(err);
  }
}

async function onLogout() {
  try {
    await api.logout();
    clearSession();
    setBanner('ok', 'Logged out');
    location.hash = 'login';
  } catch (err) {
    fail(err);
  }
}

async function onDeleteAccount() {
  if (!confirm(`Delete account ${state.user?.email}?`)) return;
  try {
    await api.deleteAccount();
    clearSession();
    setBanner('ok', 'Account deleted');
    location.hash = 'login';
  } catch (err) {
    fail(err);
  }
}

async function loadConversations() {
  state.conversations = await api.listConversations({ limit: 100, offset: 0 });
  if (!state.selectedConversationId) return;
  const stillThere = state.conversations.items.some((item) => item.id === state.selectedConversationId);
  if (!stillThere) {
    state.selectedConversationId = null;
    state.conversation = null;
    return;
  }
  state.conversation = await api.getConversation(state.selectedConversationId);
}

async function openConversation(id) {
  state.selectedConversationId = id;
  try {
    state.conversation = await api.getConversation(id);
    setBanner('', '');
    render();
  } catch (err) {
    fail(err);
  }
}

async function onCreateConversation(event) {
  event.preventDefault();
  const data = new FormData(event.target);
  const id = crypto.randomUUID();
  const title = String(data.get('title') || '').trim();
  const body = { id };
  if (title) body.title = title;
  try {
    await api.createConversation(body);
    state.selectedConversationId = id;
    await loadConversations();
    setBanner('ok', 'Conversation created');
    render();
  } catch (err) {
    fail(err);
  }
}

async function removeConversation(id) {
  if (!confirm('Delete this conversation?')) return;
  try {
    await api.deleteConversation(id);
    if (state.selectedConversationId === id) {
      state.selectedConversationId = null;
      state.conversation = null;
    }
    await loadConversations();
    setBanner('ok', 'Conversation deleted');
    render();
  } catch (err) {
    fail(err);
  }
}

async function onSend(event) {
  event.preventDefault();
  if (!state.selectedConversationId || state.sending) return;
  const data = new FormData(event.target);
  const tags = parseTags(String(data.get('tags') || ''));
  state.draft = String(data.get('content') || '');
  const body = {
    id: crypto.randomUUID(),
    content: state.draft,
    topK: Number(data.get('topK') || 5),
    historyCapacity: Number(data.get('historyCapacity') || 3),
  };
  if (tags.length) body.tags = tags;
  state.sending = true;
  setBanner('ok', 'Waiting for the assistant…');
  render();
  try {
    await api.sendMessage(state.selectedConversationId, body);
    await loadConversations();
    state.draft = '';
    setBanner('ok', 'Message sent');
  } catch (err) {
    state.sending = false;
    fail(err);
    return;
  }
  state.sending = false;
  if (state.view === 'chat') render();
}

async function loadDocuments() {
  state.documents = await api.listDocuments({
    limit: 100,
    offset: 0,
    tag: state.docFilter.tag,
    indexStatus: state.docFilter.indexStatus,
  });
}

async function openDocument(id) {
  state.selectedDocumentId = id;
  try {
    state.document = await api.getDocument(id);
    state.indexJob = null;
    setBanner('', '');
    render();
  } catch (err) {
    fail(err);
  }
}

async function onFilterDocuments(event) {
  event.preventDefault();
  const data = new FormData(event.target);
  state.docFilter = {
    tag: String(data.get('tag') || '').trim(),
    indexStatus: String(data.get('indexStatus') || ''),
  };
  try {
    await loadDocuments();
    setBanner('', '');
    render();
  } catch (err) {
    fail(err);
  }
}

async function onCreateDocument(event) {
  event.preventDefault();
  const data = new FormData(event.target);
  const id = crypto.randomUUID();
  const tagsText = String(data.get('tags') || '');
  const tags = parseTags(tagsText);
  state.docDraft = {
    title: String(data.get('title') || ''),
    tags: tagsText,
    content: String(data.get('content') || ''),
  };
  const body = {
    id,
    title: state.docDraft.title,
    content: state.docDraft.content,
  };
  if (tags.length) body.tags = tags;
  try {
    await api.createDocument(body);
    state.docDraft = { title: '', tags: '', content: '' };
    state.selectedDocumentId = id;
    await loadDocuments();
    state.document = await api.getDocument(id);
    state.indexJob = null;
    setBanner('ok', 'Document created');
    render();
  } catch (err) {
    fail(err);
  }
}

async function onUpdateDocument(event) {
  event.preventDefault();
  if (!state.selectedDocumentId) return;
  const data = new FormData(event.target);
  try {
    await api.updateDocument(state.selectedDocumentId, {
      title: String(data.get('title') || ''),
      content: String(data.get('content') || ''),
      tags: parseTags(String(data.get('tags') || '')),
    });
    await loadDocuments();
    state.document = await api.getDocument(state.selectedDocumentId);
    setBanner('ok', 'Document saved');
    render();
  } catch (err) {
    fail(err);
  }
}

async function removeDocument(id) {
  if (!confirm('Delete this document?')) return;
  try {
    await api.deleteDocument(id);
    if (state.selectedDocumentId === id) {
      state.selectedDocumentId = null;
      state.document = null;
      state.indexJob = null;
    }
    await loadDocuments();
    setBanner('ok', 'Document deleted');
    render();
  } catch (err) {
    fail(err);
  }
}

async function onReindex() {
  if (!state.selectedDocumentId) return;
  try {
    await api.reindexDocument(state.selectedDocumentId);
    state.indexJob = await api.getDocumentIndexStatus(state.selectedDocumentId);
    state.document = await api.getDocument(state.selectedDocumentId);
    await loadDocuments();
    setBanner('ok', 'Reindex accepted');
    render();
  } catch (err) {
    fail(err);
  }
}

async function onIndexStatus() {
  if (!state.selectedDocumentId) return;
  try {
    state.indexJob = await api.getDocumentIndexStatus(state.selectedDocumentId);
    setBanner('', '');
    render();
  } catch (err) {
    fail(err);
  }
}

function enter(view) {
  if ((view === 'chat' || view === 'documents') && !state.user) {
    if (currentHash() !== 'login') location.hash = 'login';
    return;
  }
  if ((view === 'login' || view === 'register') && state.user) {
    if (currentHash() !== 'chat') location.hash = 'chat';
    return;
  }
  if (view === 'documents' && state.user && !isAdmin()) {
    state.view = 'documents';
    render();
    return;
  }

  state.view = view || 'login';
  render();
  void loadView();
}

async function loadView() {
  try {
    if (state.view === 'chat') {
      await loadConversations();
      if (state.bannerKind === 'error') setBanner('', '');
      render();
    } else if (state.view === 'documents' && isAdmin()) {
      await loadDocuments();
      if (state.bannerKind === 'error') setBanner('', '');
      render();
    }
  } catch (err) {
    fail(err);
  }
}

async function boot() {
  loadSession();
  if (api.getToken()) {
    try {
      state.user = await api.getCurrentUser();
      localStorage.setItem(sessionKey, JSON.stringify({
        accessToken: api.getToken(),
        user: state.user,
      }));
    } catch (err) {
      if (err.status === 401) clearSession();
      else setBanner('error', err.message);
    }
  }

  window.addEventListener('hashchange', () => enter(currentHash()));
  const initial = currentHash() || (state.user ? 'chat' : 'login');
  if (currentHash() !== initial) location.hash = initial;
  else enter(initial);
}

boot();
