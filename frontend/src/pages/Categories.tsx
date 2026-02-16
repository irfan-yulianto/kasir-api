import { useEffect, useState } from "react"
import { api, type Category } from "@/lib/api"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
  Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle,
} from "@/components/ui/dialog"
import { toast } from "sonner"
import { Plus, Pencil, Trash2, Loader2, Tag } from "lucide-react"

export default function Categories() {
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)

  const [formOpen, setFormOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [selected, setSelected] = useState<Category | null>(null)
  const [form, setForm] = useState({ name: "", description: "" })
  const [saving, setSaving] = useState(false)

  const fetchData = async () => {
    try {
      setCategories(await api.getCategories())
    } catch {
      toast.error("Gagal memuat kategori")
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [])

  const openCreate = () => {
    setSelected(null)
    setForm({ name: "", description: "" })
    setFormOpen(true)
  }

  const openEdit = (cat: Category) => {
    setSelected(cat)
    setForm({ name: cat.name, description: cat.description })
    setFormOpen(true)
  }

  const handleSave = async () => {
    if (!form.name.trim()) { toast.error("Nama kategori wajib diisi"); return }
    setSaving(true)
    try {
      if (selected) {
        await api.updateCategory(selected.id, form)
        toast.success("Kategori berhasil diupdate")
      } else {
        await api.createCategory(form)
        toast.success("Kategori berhasil ditambahkan")
      }
      setFormOpen(false)
      fetchData()
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Gagal menyimpan kategori")
    } finally {
      setSaving(false)
    }
  }

  const openDelete = (cat: Category) => {
    setSelected(cat)
    setDeleteOpen(true)
  }

  const handleDelete = async () => {
    if (!selected) return
    setSaving(true)
    try {
      await api.deleteCategory(selected.id)
      toast.success("Kategori berhasil dihapus")
      setDeleteOpen(false)
      fetchData()
    } catch {
      toast.error("Gagal menghapus kategori")
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="flex-1 overflow-y-auto p-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Kategori</h1>
        <Button onClick={openCreate}><Plus className="h-4 w-4" /> Tambah Kategori</Button>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-20">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : categories.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-20 text-muted-foreground">
          <Tag className="h-12 w-12 mb-4" />
          <p>Belum ada kategori</p>
        </div>
      ) : (
        <div className="border rounded-lg overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-muted/50">
              <tr>
                <th className="text-left p-3 font-medium">ID</th>
                <th className="text-left p-3 font-medium">Nama</th>
                <th className="text-left p-3 font-medium">Deskripsi</th>
                <th className="text-right p-3 font-medium">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {categories.map((cat) => (
                <tr key={cat.id} className="border-t">
                  <td className="p-3 text-muted-foreground">{cat.id}</td>
                  <td className="p-3 font-medium">{cat.name}</td>
                  <td className="p-3 text-muted-foreground">{cat.description || "-"}</td>
                  <td className="p-3 text-right">
                    <div className="flex justify-end gap-1">
                      <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => openEdit(cat)}>
                        <Pencil className="h-3 w-3" />
                      </Button>
                      <Button variant="ghost" size="icon" className="h-8 w-8 text-destructive" onClick={() => openDelete(cat)}>
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
            <DialogTitle>{selected ? "Edit Kategori" : "Tambah Kategori"}</DialogTitle>
            <DialogDescription>{selected ? "Ubah informasi kategori" : "Tambahkan kategori baru"}</DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div><label className="text-sm font-medium">Nama</label><Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} /></div>
            <div><label className="text-sm font-medium">Deskripsi</label><Input value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} /></div>
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
            <DialogTitle>Hapus Kategori</DialogTitle>
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
