const API_BASE = import.meta.env.VITE_API_URL || ""

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const token = localStorage.getItem("token")
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), 30000)

  let res: Response
  try {
    res = await fetch(`${API_BASE}${path}`, {
      ...options,
      signal: controller.signal,
      headers: {
        "Content-Type": "application/json",
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...options?.headers,
      },
    })
  } catch (err) {
    clearTimeout(timeoutId)
    if (err instanceof DOMException && err.name === "AbortError") {
      throw new Error("Request timeout — server tidak merespons")
    }
    throw err
  } finally {
    clearTimeout(timeoutId)
  }

  if (res.status === 401 && !path.startsWith("/api/auth/")) {
    localStorage.removeItem("token")
    localStorage.removeItem("user")
    if (window.location.pathname !== "/login") {
      window.location.href = "/login"
    }
    throw new Error("Sesi berakhir, silakan login kembali")
  }

  if (!res.ok) {
    let message = `Request failed: ${res.status}`
    try {
      const body = await res.json()
      if (body.error) message = body.error
    } catch {
      const text = await res.text().catch(() => "")
      if (text) message = text
    }
    throw new Error(message)
  }

  return res.json()
}

// Models
export interface Product {
  id: number
  name: string
  price: number
  stock: number
  category_id?: number | null
  category?: Category
}

export interface Category {
  id: number
  name: string
  description: string
}

export interface Transaction {
  id: number
  total_amount: number
  created_at: string
  details: TransactionDetail[]
}

export interface TransactionDetail {
  id: number
  transaction_id: number
  product_id: number
  product_name?: string
  quantity: number
  subtotal: number
}

export interface CheckoutItem {
  product_id: number
  quantity: number
}

export interface CheckoutRequest {
  items: CheckoutItem[]
}

export interface SalesSummary {
  total_revenue: number
  total_transactions: number
  top_products: TopProduct[]
}

export interface TopProduct {
  product_id: number
  product_name: string
  total_sold: number
}

export interface AuthResponse {
  token: string
  user: { id: number; username: string; role: string; created_at: string }
}

// Auth API
export const authApi = {
  login(username: string, password: string): Promise<AuthResponse> {
    return request("/api/auth/login", {
      method: "POST",
      body: JSON.stringify({ username, password }),
    })
  },

  register(username: string, password: string, role: string): Promise<AuthResponse> {
    return request("/api/auth/register", {
      method: "POST",
      body: JSON.stringify({ username, password, role }),
    })
  },
}

// API
export const api = {
  // Products
  getProducts(): Promise<Product[]> {
    return request("/api/products?include_category=true")
  },

  getProduct(id: number): Promise<Product> {
    return request(`/api/products/${id}?include_category=true`)
  },

  createProduct(data: Partial<Product>): Promise<Product> {
    return request("/api/products", {
      method: "POST",
      body: JSON.stringify(data),
    })
  },

  updateProduct(id: number, data: Partial<Product>): Promise<Product> {
    return request(`/api/products/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    })
  },

  deleteProduct(id: number): Promise<{ message: string }> {
    return request(`/api/products/${id}`, {
      method: "DELETE",
    })
  },

  // Categories
  getCategories(): Promise<Category[]> {
    return request("/api/categories")
  },

  createCategory(data: Partial<Category>): Promise<Category> {
    return request("/api/categories", {
      method: "POST",
      body: JSON.stringify(data),
    })
  },

  updateCategory(id: number, data: Partial<Category>): Promise<Category> {
    return request(`/api/categories/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    })
  },

  deleteCategory(id: number): Promise<{ message: string }> {
    return request(`/api/categories/${id}`, {
      method: "DELETE",
    })
  },

  // Transactions
  checkout(data: CheckoutRequest): Promise<Transaction> {
    return request("/api/checkout", {
      method: "POST",
      body: JSON.stringify(data),
    })
  },

  getTransactions(startDate?: string, endDate?: string): Promise<Transaction[]> {
    const params = new URLSearchParams()
    if (startDate) params.set("start_date", startDate)
    if (endDate) params.set("end_date", endDate)
    const qs = params.toString()
    return request(`/api/transactions${qs ? `?${qs}` : ""}`)
  },

  // Reports
  getReport(startDate?: string, endDate?: string): Promise<SalesSummary> {
    const params = new URLSearchParams()
    if (startDate) params.set("start_date", startDate)
    if (endDate) params.set("end_date", endDate)
    const qs = params.toString()
    return request(`/api/report/today${qs ? `?${qs}` : ""}`)
  },
}
