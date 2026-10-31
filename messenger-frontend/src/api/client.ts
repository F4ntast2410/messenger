const BASE_URL = '/api/v1';

class ApiError extends Error {
  status: number;

    constructor(status: number, message: string) {
      super(message);
      this.status = status;
    }
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
    const res = await fetch(`${BASE_URL}${path}`, {
      ...options,
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    });

    if (!res.ok) {
      throw new ApiError(res.status, `Request failed: ${res.status}`);
    }

    if (res.status === 204) return undefined as T
    const text = await res.text()
    return (text ? JSON.parse(text) : undefined) as T
}

export { request, ApiError };