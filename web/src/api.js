const services = {
  user: '/svc/user/api',
  knowledge: '/svc/knowledge/api',
  conversation: '/svc/conversation/api',
};

let token = null;

export function setToken(value) {
  token = value || null;
}

export function getToken() {
  return token;
}

export class ApiError extends Error {
  constructor(message, status, slug) {
    super(message);
    this.status = status;
    this.slug = slug;
  }
}

export async function request(service, path, options = {}) {
  const url = new URL(`${services[service]}${path}`, window.location.origin);
  for (const [key, value] of Object.entries(options.query || {})) {
    if (value !== undefined && value !== null && value !== '') {
      url.searchParams.set(key, String(value));
    }
  }

  const headers = { Accept: 'application/json' };
  if (token) headers.Authorization = `Bearer ${token}`;
  if (options.body !== undefined) headers['Content-Type'] = 'application/json';

  let response;
  try {
    response = await fetch(url, {
      method: options.method || 'GET',
      headers,
      body: options.body !== undefined ? JSON.stringify(options.body) : undefined,
    });
  } catch (err) {
    throw new ApiError(`Network error calling ${service} (${err.message})`, 0, 'network');
  }

  if (response.status === 204 || response.status === 202) return null;

  const text = await response.text();
  let data = null;
  if (text) {
    try {
      data = JSON.parse(text);
    } catch {
      data = { message: text };
    }
  }

  if (!response.ok) {
    const message = data?.message || response.statusText || 'Request failed';
    throw new ApiError(message, response.status, data?.slug);
  }
  return data;
}

export const register = (body) => request('user', '/auth/register', { method: 'POST', body });
export const login = (body) => request('user', '/auth/login', { method: 'POST', body });
export const logout = () => request('user', '/auth/logout', { method: 'POST' });
export const getCurrentUser = () => request('user', '/users/me');
export const deleteAccount = () => request('user', '/users/me', { method: 'DELETE' });

export const listDocuments = (query) => request('knowledge', '/documents', { query });
export const getDocument = (id) => request('knowledge', `/documents/${id}`);
export const createDocument = (body) => request('knowledge', '/documents', { method: 'POST', body });
export const updateDocument = (id, body) => request('knowledge', `/documents/${id}`, { method: 'PUT', body });
export const deleteDocument = (id) => request('knowledge', `/documents/${id}`, { method: 'DELETE' });
export const reindexDocument = (id) => request('knowledge', `/documents/${id}/reindex`, { method: 'POST' });
export const getDocumentIndexStatus = (id) => request('knowledge', `/documents/${id}/index-status`);

export const listConversations = (query) => request('conversation', '/conversations', { query });
export const createConversation = (body) => request('conversation', '/conversations', { method: 'POST', body });
export const getConversation = (id) => request('conversation', `/conversations/${id}`);
export const deleteConversation = (id) => request('conversation', `/conversations/${id}`, { method: 'DELETE' });
export const sendMessage = (conversationId, body) =>
  request('conversation', `/conversations/${conversationId}/messages`, { method: 'POST', body });
