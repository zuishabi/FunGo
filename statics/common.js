// file: `statics/common.js`
console.debug('[common.js] loaded');

const TOKEN_KEY = 'token';
function getAuthToken(){ return localStorage.getItem(TOKEN_KEY) || ''; }
function setAuthToken(token){ if(token) localStorage.setItem(TOKEN_KEY, token); else localStorage.removeItem(TOKEN_KEY); }

async function apiFetch(url, options = {}) {
    const opts = { method: 'GET', headers: {}, credentials: 'include', ...options };
    const token = getAuthToken();
    if (token) opts.headers.Authorization = `Bearer ${token}`;
    if (!(opts.body instanceof FormData) && opts.body && typeof opts.body === 'object') {
        opts.headers['Content-Type'] = 'application/json';
        opts.body = JSON.stringify(opts.body);
    }
    const resp = await fetch(url, opts);
    const text = await resp.text();
    let json;
    try { json = text ? JSON.parse(text) : null; } catch (e) {
        if (!resp.ok) throw new Error(text || `请求失败: ${resp.status}`);
        return text;
    }
    if (!resp.ok) throw new Error(json?.msg || text || `请求失败: ${resp.status}`);
    if (typeof json?.code !== 'undefined' && json.code !== 0) throw new Error(json?.msg || text || '请求返回错误');
    return json?.data ?? json;
}

async function doLogin(username, password) {
    const data = await apiFetch('/api/user/login', { method: 'POST', body: { user_name: username, password } });
    const token = data?.token || data?.data?.token;
    if (!token) throw new Error('登录返回缺少 token');
    setAuthToken(token);
    await updateAuthUI();
    return data;
}

async function doRegister(username, password) {
    if (!username || !password) throw new Error('用户名和密码不能为空');
    await apiFetch('/api/user/register', { method: 'POST', body: { user_name: username, password } });
    return true;
}

function logout() {
    setAuthToken('');
    updateAuthUI();
}

async function loadHeader(containerSelector = '#header-container') {
    const el = document.querySelector(containerSelector);
    if (!el) return;
    const resp = await fetch('/header');
    el.innerHTML = await resp.text();
    await updateAuthUI();
}

async function updateAuthUI() {
    const token = getAuthToken();
    const loginLink = document.getElementById('login-link');
    const registerLink = document.getElementById('register-link');
    const avatarWrap = document.getElementById('avatar-wrap');
    const avatar = document.getElementById('avatar');

    if (!loginLink) return;
    if (!token) {
        loginLink.style.display = '';
        if (registerLink) registerLink.style.display = '';
        if (avatarWrap) avatarWrap.style.display = 'none';
        return;
    }

    try {
        const info = await apiFetch('/api/user/selfInfo');
        const name = info?.user_name || info?.UserName || '用户';
        const initial = (name.slice(-1) || 'U').toUpperCase();
        if (avatar) avatar.textContent = initial;
        loginLink.style.display = 'none';
        if (registerLink) registerLink.style.display = 'none';
        if (avatarWrap) avatarWrap.style.display = '';
    } catch {
        setAuthToken('');
        loginLink.style.display = '';
        if (registerLink) registerLink.style.display = '';
        if (avatarWrap) avatarWrap.style.display = 'none';
    }
}

/* 全局代理：处理 avatar 菜单与退出 */
if (!window.__header_simple_delegate_bound) {
    window.__header_simple_delegate_bound = true;
    document.addEventListener('click', async (e) => {
        const raw = e.target;
        const el = (raw && raw.nodeType === 3) ? raw.parentElement : raw;

        // avatar 点击显示/隐藏菜单
        if (el && el.closest && el.closest('#avatar')) {
            e.stopPropagation();
            const m = document.getElementById('avatar-menu');
            if (m) m.style.display = m.style.display === 'block' ? 'none' : 'block';
            return;
        }

        // menu logout
        if (el && el.closest && el.closest('#menu-logout')) {
            e.preventDefault();
            logout();
            await loadHeader('#header-container');
            return;
        }

        // 点击其他地方关闭 avatar 菜单
        if (!el || !el.closest || !el.closest('#avatar-wrap')) {
            const m = document.getElementById('avatar-menu');
            if (m) m.style.display = 'none';
        }

        // 导航按钮处理（保留原有 nav id）
        if (el && el.closest && el.closest('#nav-posts')) { location.href = '/'; return; }
        // 直接跳到标准 section 路由，避免中间 /hot 导致状态不一致
        if (el && el.closest && el.closest('#nav-game')) { location.href = '/gameList'; return; }
        if (el && el.closest && el.closest('#nav-live')) { location.href = '/liveList'; return; }
        if (el && el.closest && el.closest('#nav-community')) { location.hash = '/community';  }
    });
}

// 导出接口供登录/注册页面使用
window.apiFetch = apiFetch;
window.loadHeader = loadHeader;
window.updateAuthUI = updateAuthUI;
window.doLogin = doLogin;
window.doRegister = doRegister;
window.logout = logout;
