import { useState } from "react"
import { api, type Product, type Transaction } from "@/lib/api"
import { formatRupiah } from "@/lib/utils"
import { Toaster } from "@/components/ui/sonner"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Separator } from "@/components/ui/separator"
import {
  Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle,
} from "@/components/ui/dialog"
import { toast } from "sonner"
import {
  Eye, Pencil, Trash2, ShoppingCart, Package, Loader2,
  Minus, Plus, CheckCircle2,
} from "lucide-react"

interface CartItem {
  product: Product
  quantity: number
}

export default function App() {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [cart, setCart] = useState<CartItem[]>([])
  const [receipt, setReceipt] = useState<Transaction | null>(null)
  const [checkingOut, setCheckingOut] = useState(false)

  // Dialog states
  const [selected, setSelected] = useState<Product | null>(null)
  const [detailOpen, setDetailOpen] = useState(false)
  const [editOpen, setEditOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [editForm, setEditForm] = useState({ name: "", price: 0, stock: 0 })
  const [saving, setSaving] = useState(false)

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

  useState(() => { fetchProducts() })

  // Cart actions
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

  // Product actions
  const handleDetail = async (id: number) => {
    try {
      const product = await api.getProduct(id)
      setSelected(product)
      setDetailOpen(true)
    } catch {
      toast.error("Gagal memuat detail produk")
    }
  }

  const openEdit = (product: Product) => {
    setSelected(product)
    setEditForm({ name: product.name, price: product.price, stock: product.stock })
    setEditOpen(true)
  }

  const handleUpdate = async () => {
    if (!selected) return
    setSaving(true)
    try {
      await api.updateProduct(selected.id, editForm)
      toast.success("Produk berhasil diupdate")
      setEditOpen(false)
      fetchProducts()
    } catch {
      toast.error("Gagal mengupdate produk")
    } finally {
      setSaving(false)
    }
  }

  const openDelete = (product: Product) => {
    setSelected(product)
    setDeleteOpen(true)
  }

  const handleDelete = async () => {
    if (!selected) return
    setSaving(true)
    try {
      await api.deleteProduct(selected.id)
      toast.success("Produk berhasil dihapus")
      setDeleteOpen(false)
      fetchProducts()
    } catch {
      toast.error("Gagal menghapus produk")
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="sticky top-0 z-40 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-14 items-center justify-between px-6">
          <h1 className="text-lg font-bold">Kasir App</h1>
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <ShoppingCart className="h-4 w-4" />
            <span>{cartCount} item</span>
            <span className="font-semibold text-foreground">{formatRupiah(cartTotal)}</span>
          </div>
        </div>
      </header>

      {/* Main: Products left, Cart right */}
      <div className="flex h-[calc(100vh-3.5rem)]">
        {/* Product List */}
        <main className="flex-1 overflow-y-auto p-6">
          {loading ? (
            <div className="flex items-center justify-center py-20">
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
          ) : products.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-20 text-muted-foreground">
              <Package className="h-12 w-12 mb-4" />
              <p>Belum ada produk</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-3">
              {products.map((product) => (
                <Card key={product.id} className="flex flex-col">
                  <CardHeader className="pb-3">
                    <div className="flex items-start justify-between">
                      <CardTitle className="text-base">{product.name}</CardTitle>
                      <Badge variant={product.stock > 0 ? "secondary" : "destructive"}>
                        {product.stock > 0 ? `Stok: ${product.stock}` : "Habis"}
                      </Badge>
                    </div>
                  </CardHeader>
                  <CardContent className="flex-1">
                    <p className="text-2xl font-bold">{formatRupiah(product.price)}</p>
                  </CardContent>
                  <CardFooter className="gap-2 flex-wrap">
                    <Button variant="outline" size="sm" onClick={() => handleDetail(product.id)}>
                      <Eye className="h-3 w-3" />
                    </Button>
                    <Button variant="outline" size="sm" onClick={() => openEdit(product)}>
                      <Pencil className="h-3 w-3" />
                    </Button>
                    <Button variant="outline" size="sm" onClick={() => openDelete(product)}>
                      <Trash2 className="h-3 w-3" />
                    </Button>
                    <Button size="sm" className="ml-auto" disabled={product.stock <= 0} onClick={() => addToCart(product)}>
                      <ShoppingCart className="h-3 w-3" /> Tambah
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
            /* Receipt */
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
            /* Empty cart */
            <div className="flex-1 flex flex-col items-center justify-center text-muted-foreground p-4">
              <ShoppingCart className="h-10 w-10 mb-3" />
              <p className="text-sm">Keranjang kosong</p>
              <p className="text-xs">Klik "Tambah" pada produk</p>
            </div>
          ) : (
            /* Cart items */
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
                      <button className="text-xs text-destructive hover:underline" onClick={() => removeFromCart(item.product.id)}>Hapus</button>
                    </div>
                  </div>
                ))}
              </div>

              {/* Checkout footer */}
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

      {/* Detail Dialog */}
      <Dialog open={detailOpen} onOpenChange={setDetailOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Detail Produk</DialogTitle>
            <DialogDescription>Informasi lengkap produk</DialogDescription>
          </DialogHeader>
          {selected && (
            <div className="space-y-3">
              <div className="flex justify-between"><span className="text-muted-foreground">Nama</span><span className="font-medium">{selected.name}</span></div>
              <div className="flex justify-between"><span className="text-muted-foreground">Harga</span><span className="font-medium">{formatRupiah(selected.price)}</span></div>
              <div className="flex justify-between"><span className="text-muted-foreground">Stok</span><span className="font-medium">{selected.stock}</span></div>
              {selected.category_name && (
                <div className="flex justify-between"><span className="text-muted-foreground">Kategori</span><Badge variant="outline">{selected.category_name}</Badge></div>
              )}
            </div>
          )}
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      <Dialog open={editOpen} onOpenChange={setEditOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Produk</DialogTitle>
            <DialogDescription>Ubah informasi produk</DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div><label className="text-sm font-medium">Nama</label><Input value={editForm.name} onChange={(e) => setEditForm({ ...editForm, name: e.target.value })} /></div>
            <div><label className="text-sm font-medium">Harga</label><Input type="number" value={editForm.price} onChange={(e) => setEditForm({ ...editForm, price: Number(e.target.value) })} /></div>
            <div><label className="text-sm font-medium">Stok</label><Input type="number" value={editForm.stock} onChange={(e) => setEditForm({ ...editForm, stock: Number(e.target.value) })} /></div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setEditOpen(false)}>Batal</Button>
            <Button onClick={handleUpdate} disabled={saving}>{saving && <Loader2 className="h-4 w-4 animate-spin" />} Simpan</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Dialog */}
      <Dialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Hapus Produk</DialogTitle>
            <DialogDescription>Apakah kamu yakin ingin menghapus <strong>{selected?.name}</strong>? Aksi ini tidak bisa dibatalkan.</DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteOpen(false)}>Batal</Button>
            <Button variant="destructive" onClick={handleDelete} disabled={saving}>{saving && <Loader2 className="h-4 w-4 animate-spin" />} Hapus</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Toaster richColors position="top-right" />
    </div>
  )
}
