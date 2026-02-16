const API_BASE = import.meta.env.VITE_API_URL || ""
const API_KEY = import.meta.env.VITE_API_KEY || "your-secret-api-key-here"

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
  })

  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `Request failed: ${res.status}`)
  }

  return res.json()
}

function authHeaders(): HeadersInit {
  return { "X-API-Key": API_KEY }
}

// Models
export interface Product {
  id: number
  name: string
  price: number
  stock: number
  category_id?: number | null
  category_name?: string
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

// API
export const api = {
  // Products
  getProducts(): Promise<Product[]> {
    return request("/api/products?include_category=true")
  },

  getProduct(id: number): Promise<Product> {
    return request(`/api/products/${id}?include_category=true`, { headers: authHeaders() })
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
      headers: authHeaders(),
      body: JSON.stringify(data),
    })
  },

  deleteProduct(id: number): Promise<{ message: string }> {
    return request(`/api/products/${id}`, {
      method: "DELETE",
      headers: authHeaders(),
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
      headers: authHeaders(),
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
