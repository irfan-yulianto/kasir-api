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

export interface Product {
  id: number
  name: string
  price: number
  stock: number
  category_id?: number | null
  category_name?: string
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

export const api = {
  getProducts(): Promise<Product[]> {
    return request("/api/products")
  },

  getProduct(id: number): Promise<Product> {
    return request(`/api/products/${id}`, { headers: authHeaders() })
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

  checkout(data: CheckoutRequest): Promise<Transaction> {
    return request("/api/checkout", {
      method: "POST",
      headers: authHeaders(),
      body: JSON.stringify(data),
    })
  },
}
