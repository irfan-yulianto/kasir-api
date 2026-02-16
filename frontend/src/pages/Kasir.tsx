import { useEffect, useState } from "react"
import { api, type Product, type Transaction } from "@/lib/api"
import { formatRupiah } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Separator } from "@/components/ui/separator"
import { toast } from "sonner"
import {
  ShoppingCart, Package, Loader2, Minus, Plus, CheckCircle2, Search,
} from "lucide-react"

interface CartItem {
  product: Product
  quantity: number
}

export default function Kasir() {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState("")
  const [cart, setCart] = useState<CartItem[]>([])
  const [receipt, setReceipt] = useState<Transaction | null>(null)
  const [checkingOut, setCheckingOut] = useState(false)

  const cartTotal = cart.reduce((sum, item) => sum + item.product.price * item.quantity, 0)
  const cartCount = cart.reduce((sum, item) => sum + item.quantity, 0)

  const fetchProducts = async () => {
    try {
      const data = await api.getProducts()
      setProducts(data)
    } catch {
      toast.error("Gagal memuat produk")
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchProducts() }, [])

  const filtered = products.filter((p) =>
    p.name.toLowerCase().includes(search.toLowerCase())
  )

  const addToCart = (product: Product) => {
    setCart((prev) => {
      const existing = prev.find((item) => item.product.id === product.id)
      if (existing) {
        if (existing.quantity >= product.stock) return prev
        return prev.map((item) =>
          item.product.id === product.id ? { ...item, quantity: item.quantity + 1 } : item
        )
      }
      return [...prev, { product, quantity: 1 }]
    })
    toast.success(`${product.name} ditambahkan`)
  }

  const updateQuantity = (productId: number, qty: number) => {
    if (qty <= 0) {
      setCart((prev) => prev.filter((item) => item.product.id !== productId))
      return
    }
    setCart((prev) =>
      prev.map((item) => (item.product.id === productId ? { ...item, quantity: qty } : item))
    )
  }

  const removeFromCart = (productId: number) => {
    setCart((prev) => prev.filter((item) => item.product.id !== productId))
  }

  const handleCheckout = async () => {
    if (cart.length === 0) return
    setCheckingOut(true)
    try {
      const result = await api.checkout({
        items: cart.map((item) => ({ product_id: item.product.id, quantity: item.quantity })),
      })
      setReceipt(result)
      setCart([])
      fetchProducts()
      toast.success("Checkout berhasil!")
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Checkout gagal")
    } finally {
      setCheckingOut(false)
    }
  }

  return (
    <div className="flex h-full">
      {/* Product List */}
      <main className="flex-1 overflow-y-auto p-6">
        <div className="mb-4">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Cari produk..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-10"
            />
          </div>
        </div>

        {loading ? (
          <div className="flex items-center justify-center py-20">
            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
          </div>
        ) : filtered.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-20 text-muted-foreground">
            <Package className="h-12 w-12 mb-4" />
            <p>{search ? "Produk tidak ditemukan" : "Belum ada produk"}</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-3">
            {filtered.map((product) => (
              <Card key={product.id} className="flex flex-col">
                <CardHeader className="pb-3">
                  <div className="flex items-start justify-between">
                    <CardTitle className="text-base">{product.name}</CardTitle>
                    <Badge variant={product.stock > 0 ? "secondary" : "destructive"}>
                      {product.stock > 0 ? `Stok: ${product.stock}` : "Habis"}
                    </Badge>
                  </div>
                  {product.category?.name && (
                    <Badge variant="outline" className="w-fit">{product.category.name}</Badge>
                  )}
                </CardHeader>
                <CardContent className="flex-1">
                  <p className="text-2xl font-bold">{formatRupiah(product.price)}</p>
                </CardContent>
                <CardFooter>
                  <Button className="w-full" disabled={product.stock <= 0} onClick={() => addToCart(product)}>
                    <ShoppingCart className="h-4 w-4" /> Tambah ke Keranjang
                  </Button>
                </CardFooter>
              </Card>
            ))}
          </div>
        )}
      </main>

      {/* Cart Sidebar */}
      <aside className="w-96 border-l bg-card flex flex-col overflow-hidden">
        <div className="p-4 border-b">
          <h2 className="font-semibold flex items-center gap-2">
            <ShoppingCart className="h-4 w-4" /> Keranjang
            {cartCount > 0 && <Badge variant="secondary">{cartCount}</Badge>}
          </h2>
        </div>

        {receipt ? (
          <div className="flex-1 overflow-y-auto p-4">
            <div className="text-center mb-4">
              <CheckCircle2 className="h-10 w-10 text-green-600 mx-auto mb-2" />
              <p className="font-semibold">Transaksi Berhasil!</p>
              <p className="text-sm text-muted-foreground">#{receipt.id} - {new Date(receipt.created_at).toLocaleString("id-ID")}</p>
            </div>
            <Separator className="my-3" />
            {receipt.details.map((d) => (
              <div key={d.id} className="flex justify-between text-sm py-1">
                <span>{d.product_name} x{d.quantity}</span>
                <span>{formatRupiah(d.subtotal)}</span>
              </div>
            ))}
            <Separator className="my-3" />
            <div className="flex justify-between font-bold">
              <span>Total</span>
              <span>{formatRupiah(receipt.total_amount)}</span>
            </div>
            <Button className="w-full mt-4" onClick={() => setReceipt(null)}>Transaksi Baru</Button>
          </div>
        ) : cart.length === 0 ? (
          <div className="flex-1 flex flex-col items-center justify-center text-muted-foreground p-4">
            <ShoppingCart className="h-10 w-10 mb-3" />
            <p className="text-sm">Keranjang kosong</p>
          </div>
        ) : (
          <>
            <div className="flex-1 overflow-y-auto p-4 space-y-3">
              {cart.map((item) => (
                <div key={item.product.id} className="flex items-center gap-3 rounded-lg border p-3">
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium truncate">{item.product.name}</p>
                    <p className="text-xs text-muted-foreground">{formatRupiah(item.product.price)}</p>
                  </div>
                  <div className="flex items-center gap-1">
                    <Button variant="outline" size="icon" className="h-7 w-7"
                      onClick={() => updateQuantity(item.product.id, item.quantity - 1)}>
                      <Minus className="h-3 w-3" />
                    </Button>
                    <span className="w-8 text-center text-sm font-medium">{item.quantity}</span>
                    <Button variant="outline" size="icon" className="h-7 w-7"
                      disabled={item.quantity >= item.product.stock}
                      onClick={() => updateQuantity(item.product.id, item.quantity + 1)}>
                      <Plus className="h-3 w-3" />
                    </Button>
                  </div>
                  <div className="text-right">
                    <p className="text-sm font-medium w-20 text-right">{formatRupiah(item.product.price * item.quantity)}</p>
                    <Button variant="link" size="sm" className="h-auto p-0 text-xs text-destructive" onClick={() => removeFromCart(item.product.id)}>Hapus</Button>
                  </div>
                </div>
              ))}
            </div>
            <div className="border-t p-4 space-y-3">
              <div className="flex justify-between font-bold text-lg">
                <span>Total</span>
                <span>{formatRupiah(cartTotal)}</span>
              </div>
              <Button className="w-full" size="lg" onClick={handleCheckout} disabled={checkingOut}>
                {checkingOut ? <Loader2 className="h-4 w-4 animate-spin" /> : <ShoppingCart className="h-4 w-4" />}
                {checkingOut ? "Memproses..." : "Checkout"}
              </Button>
            </div>
          </>
        )}
      </aside>
    </div>
  )
}
