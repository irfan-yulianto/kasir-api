import { useEffect, useState } from "react"
import { api, type Product, type Category } from "@/lib/api"
import { formatRupiah } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import {
  Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle,
} from "@/components/ui/dialog"
import { toast } from "sonner"
import { Plus, Pencil, Trash2, Search, Loader2, Package } from "lucide-react"

export default function Products() {
  const [products, setProducts] = useState<Product[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState("")
  const [filterCategory, setFilterCategory] = useState<number | null>(null)

  const [formOpen, setFormOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [selected, setSelected] = useState<Product | null>(null)
  const [form, setForm] = useState({ name: "", price: 0, stock: 0, category_id: 0 })
  const [saving, setSaving] = useState(false)

  const fetchData = async () => {
    try {
      const [p, c] = await Promise.all([api.getProducts(), api.getCategories()])
      setProducts(p)
      setCategories(c)
    } catch {
      toast.error("Gagal memuat data")
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [])

  const filtered = products.filter((p) => {
    const matchSearch = p.name.toLowerCase().includes(search.toLowerCase())
    const matchCategory = !filterCategory || p.category_id === filterCategory
    return matchSearch && matchCategory
  })

  const openCreate = () => {
    setSelected(null)
    setForm({ name: "", price: 0, stock: 0, category_id: 0 })
    setFormOpen(true)
  }

  const openEdit = (product: Product) => {
    setSelected(product)
    setForm({ name: product.name, price: product.price, stock: product.stock, category_id: product.category_id || 0 })
    setFormOpen(true)
  }

  const handleSave = async () => {
    if (!form.name.trim()) { toast.error("Nama produk wajib diisi"); return }
    if (form.price < 0) { toast.error("Harga tidak boleh negatif"); return }
    if (form.stock < 0) { toast.error("Stok tidak boleh negatif"); return }

    setSaving(true)
    try {
      const data = {
        name: form.name,
        price: form.price,
        stock: form.stock,
        category_id: form.category_id || undefined,
      }

      if (selected) {
        await api.updateProduct(selected.id, data)
        toast.success("Produk berhasil diupdate")
      } else {
        await api.createProduct(data)
        toast.success("Produk berhasil ditambahkan")
      }
      setFormOpen(false)
      fetchData()
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Gagal menyimpan produk")
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
      fetchData()
    } catch {
      toast.error("Gagal menghapus produk")
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="flex-1 overflow-y-auto p-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Produk</h1>
        <Button onClick={openCreate}><Plus className="h-4 w-4" /> Tambah Produk</Button>
      </div>

      <div className="flex gap-3 mb-4">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input placeholder="Cari produk..." value={search} onChange={(e) => setSearch(e.target.value)} className="pl-10" />
        </div>
        <select
          className="border rounded-md px-3 py-2 text-sm bg-background"
          value={filterCategory || ""}
          onChange={(e) => setFilterCategory(e.target.value ? Number(e.target.value) : null)}
        >
          <option value="">Semua Kategori</option>
          {categories.map((c) => (
            <option key={c.id} value={c.id}>{c.name}</option>
          ))}
        </select>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-20">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : filtered.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-20 text-muted-foreground">
          <Package className="h-12 w-12 mb-4" />
          <p>{search || filterCategory ? "Tidak ada produk yang cocok" : "Belum ada produk"}</p>
        </div>
      ) : (
        <div className="border rounded-lg overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-muted/50">
              <tr>
                <th className="text-left p-3 font-medium">Nama</th>
                <th className="text-left p-3 font-medium">Kategori</th>
                <th className="text-right p-3 font-medium">Harga</th>
                <th className="text-right p-3 font-medium">Stok</th>
                <th className="text-right p-3 font-medium">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {filtered.map((product) => (
                <tr key={product.id} className="border-t">
                  <td className="p-3 font-medium">{product.name}</td>
                  <td className="p-3">
                    {product.category_name ? <Badge variant="outline">{product.category_name}</Badge> : <span className="text-muted-foreground">-</span>}
                  </td>
                  <td className="p-3 text-right">{formatRupiah(product.price)}</td>
                  <td className="p-3 text-right">
                    <Badge variant={product.stock > 5 ? "secondary" : product.stock > 0 ? "outline" : "destructive"}>
                      {product.stock}
                    </Badge>
                  </td>
                  <td className="p-3 text-right">
                    <div className="flex justify-end gap-1">
                      <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => openEdit(product)}>
                        <Pencil className="h-3 w-3" />
                      </Button>
                      <Button variant="ghost" size="icon" className="h-8 w-8 text-destructive" onClick={() => openDelete(product)}>
                        <Trash2 className="h-3 w-3" />
                      </Button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Create/Edit Dialog */}
      <Dialog open={formOpen} onOpenChange={setFormOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{selected ? "Edit Produk" : "Tambah Produk"}</DialogTitle>
            <DialogDescription>{selected ? "Ubah informasi produk" : "Tambahkan produk baru"}</DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div><label className="text-sm font-medium">Nama</label><Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} /></div>
            <div><label className="text-sm font-medium">Harga</label><Input type="number" min={0} value={form.price} onChange={(e) => setForm({ ...form, price: Number(e.target.value) })} /></div>
            <div><label className="text-sm font-medium">Stok</label><Input type="number" min={0} value={form.stock} onChange={(e) => setForm({ ...form, stock: Number(e.target.value) })} /></div>
            <div>
              <label className="text-sm font-medium">Kategori</label>
              <select
                className="w-full border rounded-md px-3 py-2 text-sm bg-background mt-1"
                value={form.category_id}
                onChange={(e) => setForm({ ...form, category_id: Number(e.target.value) })}
              >
                <option value={0}>Tanpa Kategori</option>
                {categories.map((c) => (
                  <option key={c.id} value={c.id}>{c.name}</option>
                ))}
              </select>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setFormOpen(false)}>Batal</Button>
            <Button onClick={handleSave} disabled={saving}>{saving && <Loader2 className="h-4 w-4 animate-spin" />} Simpan</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Dialog */}
      <Dialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Hapus Produk</DialogTitle>
            <DialogDescription>Apakah kamu yakin ingin menghapus <strong>{selected?.name}</strong>?</DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteOpen(false)}>Batal</Button>
            <Button variant="destructive" onClick={handleDelete} disabled={saving}>{saving && <Loader2 className="h-4 w-4 animate-spin" />} Hapus</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
