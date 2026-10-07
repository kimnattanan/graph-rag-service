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
    links.push('<a href="#chat">Chat</a>');
    if (isAdmin()) links.push('<a href="#documents">Documents</a>');
  } else {
    links.push('<a href="#login">Login</a>');
    links.push('<a href="#register">Register</a>');
  }
  navEl.innerHTML = links.join('');

  if (state.user) {
    whoEl.innerHTML = `
      <span>${esc(state.user.username)} (${esc(state.user.role)})</span>
      <button type="button" id="logout">Logout</button>
      <button type="button" id="delete-account">Delete account</button>
    `;
  } else {
    whoEl.innerHTML = '';
  }

  if (state.banner) {
    bannerEl.hidden = false;
    bannerEl.className = state.bannerKind;
    bannerEl.textContent = state.banner;
  } else {
    bannerEl.hidden = true;
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
    <form class="auth" id="login-form">
      <h2>Login</h2>
      <label>Email <input name="email" type="email" required autocomplete="username"></label>
      <label>Password <input name="password" type="password" required autocomplete="current-password"></label>
      <button>Login</button>
      <p>Proxied to user :3001, knowledge :3000, conversation :3002.</p>
    </form>
  `;
}

function registerView() {
  return `
    <form class="auth" id="register-form">
      <h2>Register</h2>
      <label>Email <input name="email" type="email" required></label>
      <label>Username <input name="username" required maxlength="64"></label>
      <label>Password <input name="password" type="password" required minlength="8" maxlength="72"></label>
      <button>Register</button>
    </form>
  `;
}

function chatView() {
  const items = state.conversations?.items || [];
  const list = state.conversations
    ? items.map((item) => `
        <div>
          <button type="button" data-open-conversation="${esc(item.id)}" ${item.id === state.selectedConversationId ? 'disabled' : ''}>
            ${esc(item.title || '(untitled)')} (${item.messageCount})
          </button>
          <button type="button" data-delete-conversation="${esc(item.id)}">Delete</button>
        </div>
      `).join('') || '<p>No conversations.</p>'
    : `<p>${state.banner ? 'Could not load conversations.' : 'Loading…'}</p>`;

  const messages = state.conversation?.messages?.map((message) => {
    const sources = (message.sources || []).map((source) => `
      <li>${esc(source.documentTitle)} (${source.score}) — ${esc(source.text)}</li>
    `).join('');
    return `
      <article class="message">
        <div class="role">${esc(message.role)}</div>
        <div>${esc(message.content)}</div>
        ${sources ? `<details><summary>Sources (${message.sources.length})</summary><ul>${sources}</ul></details>` : ''}
      </article>
    `;
  }).join('') || (state.conversation ? '<p>No messages yet.</p>' : '<p>Select or create a conversation.</p>');

  return `
    <div class="chat">
      <aside class="side">
        <form id="new-conversation">
          <label>Title <input name="title" type="text" placeholder="optional"></label>
          <button>New conversation</button>
        </form>
        <p>${state.conversations ? `${state.conversations.total} total` : ''}</p>
        ${list}
      </aside>
      <section class="thread">
        <h2>${esc(state.conversation?.title || 'Chat')}</h2>
        <div class="messages">${messages}</div>
        <form id="composer">
          <label>Message <textarea name="content" required ${state.sending ? 'disabled' : ''}>${esc(state.draft)}</textarea></label>
          <div class="inline">
            <label>topK <input name="topK" type="number" min="1" max="50" value="5"></label>
            <label>history <input name="historyCapacity" type="number" min="0" max="100" value="3"></label>
            <label>tags <input name="tags" type="text" placeholder="a, b"></label>
          </div>
          <button ${state.sending || !state.conversation ? 'disabled' : ''}>${state.sending ? 'Waiting…' : 'Send'}</button>
        </form>
      </section>
    </div>
  `;
}

function documentsView() {
  if (!isAdmin()) {
    return '<p>Documents require an admin account (knowledge:write).</p>';
  }

  const items = state.documents?.items || [];
  const rows = items.map((item) => `
    <tr class="${item.id === state.selectedDocumentId ? 'selected' : ''}">
      <td>${esc(item.title)}</td>
      <td>${esc((item.tags || []).join(', '))}</td>
      <td>${esc(item.indexStatus)}</td>
      <td>${esc(item.updatedAt)}</td>
      <td>
        <button type="button" data-open-document="${esc(item.id)}">Open</button>
        <button type="button" data-delete-document="${esc(item.id)}">Delete</button>
      </td>
    </tr>
  `).join('');

  const doc = state.document;
  const job = state.indexJob;
  const editor = doc ? `
    <form id="edit-document">
      <h2>Edit</h2>
      <p>id ${esc(doc.id)} · index ${esc(doc.indexStatus)}</p>
      <label>Title <input name="title" type="text" required value="${esc(doc.title)}"></label>
      <label>Tags <input name="tags" type="text" value="${esc((doc.tags || []).join(', '))}"></label>
      <label>Content <textarea name="content" required>${esc(doc.content)}</textarea></label>
      <div class="actions">
        <button>Save</button>
        <button type="button" id="reindex">Reindex</button>
        <button type="button" id="index-status">Index status</button>
        <button type="button" id="delete-document">Delete</button>
      </div>
    </form>
    ${job ? `<pre>status ${esc(job.status)}
started ${esc(job.startedAt || '')}
finished ${esc(job.finishedAt || '')}
error ${esc(job.errorMessage || '')}</pre>` : ''}
  ` : '<p>Open a document to edit it.</p>';

  return `
    <h2>Documents ${state.documents ? `(${state.documents.total})` : ''}</h2>
    <form id="doc-filter" class="inline">
      <label>Tag <input name="tag" type="text" value="${esc(state.docFilter.tag)}"></label>
      <label>Index status
        <select name="indexStatus">
          ${['', 'pending', 'indexing', 'completed', 'failed'].map((status) => `
            <option value="${status}" ${state.docFilter.indexStatus === status ? 'selected' : ''}>${status || 'any'}</option>
          `).join('')}
        </select>
      </label>
      <button>Refresh</button>
    </form>
    <form id="create-document">
      <h3>Create</h3>
      <label>Title <input name="title" type="text" required value="${esc(state.docDraft.title)}"></label>
      <label>Tags <input name="tags" type="text" placeholder="comma-separated" value="${esc(state.docDraft.tags)}"></label>
      <label>Content <textarea name="content" required placeholder="markdown">${esc(state.docDraft.content)}</textarea></label>
      <button>Create</button>
    </form>
    <div class="docs">
      <div class="doc-list">
        ${state.documents ? `
          <table>
            <thead><tr><th>Title</th><th>Tags</th><th>Index</th><th>Updated</th><th></th></tr></thead>
            <tbody>${rows || '<tr><td colspan="5">No documents.</td></tr>'}</tbody>
          </table>
        ` : `<p>${state.banner ? 'Could not load documents.' : 'Loading…'}</p>`}
      </div>
      <div class="editor">${editor}</div>
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
