import { useEffect, useState } from "react"
import { format, startOfMonth } from "date-fns"
import type { DateRange } from "react-day-picker"
import { api, type Transaction } from "@/lib/api"
import { formatRupiah } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { DateRangePicker } from "@/components/ui/date-range-picker"
import { Separator } from "@/components/ui/separator"
import {
  Table, TableBody, TableCell, TableHead, TableHeader, TableRow,
} from "@/components/ui/table"
import {
  Dialog, DialogContent, DialogHeader, DialogTitle,
} from "@/components/ui/dialog"
import { toast } from "sonner"
import { Loader2, Receipt, Eye } from "lucide-react"

export default function Transactions() {
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [loading, setLoading] = useState(true)
  const [dateRange, setDateRange] = useState<DateRange | undefined>({
    from: startOfMonth(new Date()),
    to: new Date(),
  })
  const [selected, setSelected] = useState<Transaction | null>(null)
  const [detailOpen, setDetailOpen] = useState(false)

  const fetchData = async () => {
    setLoading(true)
    try {
      const data = await api.getTransactions(
        dateRange?.from ? format(dateRange.from, "yyyy-MM-dd") : undefined,
        dateRange?.to ? format(dateRange.to, "yyyy-MM-dd") : undefined,
      )
      setTransactions(data)
    } catch {
      toast.error("Gagal memuat transaksi")
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [dateRange])

  const openDetail = (t: Transaction) => {
    setSelected(t)
    setDetailOpen(true)
  }

  return (
    <div className="flex-1 overflow-y-auto p-6">
      <h1 className="text-2xl font-bold mb-6">Riwayat Transaksi</h1>

      <div className="flex gap-3 mb-4 items-center">
        <DateRangePicker value={dateRange} onChange={setDateRange} />
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-20">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : transactions.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-20 text-muted-foreground">
          <Receipt className="h-12 w-12 mb-4" />
          <p>Belum ada transaksi</p>
        </div>
      ) : (
        <div className="border rounded-lg overflow-hidden">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>ID</TableHead>
                <TableHead>Tanggal</TableHead>
                <TableHead className="text-right">Jumlah Item</TableHead>
                <TableHead className="text-right">Total</TableHead>
                <TableHead className="text-right">Aksi</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {transactions.map((t) => (
                <TableRow key={t.id}>
                  <TableCell><Badge variant="outline">#{t.id}</Badge></TableCell>
                  <TableCell>{new Date(t.created_at).toLocaleString("id-ID")}</TableCell>
                  <TableCell className="text-right">{t.details?.length || 0} item</TableCell>
                  <TableCell className="text-right font-medium">{formatRupiah(t.total_amount)}</TableCell>
                  <TableCell className="text-right">
                    <Button variant="ghost" size="sm" onClick={() => openDetail(t)}>
                      <Eye className="h-3 w-3" /> Detail
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}

      {/* Detail Dialog */}
      <Dialog open={detailOpen} onOpenChange={setDetailOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Detail Transaksi #{selected?.id}</DialogTitle>
          </DialogHeader>
          {selected && (
            <div className="space-y-3">
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Tanggal</span>
                <span>{new Date(selected.created_at).toLocaleString("id-ID")}</span>
              </div>
              <Separator />
              {selected.details?.map((d) => (
                <div key={d.id} className="flex justify-between text-sm">
                  <span>{d.product_name} x{d.quantity}</span>
                  <span>{formatRupiah(d.subtotal)}</span>
                </div>
              ))}
              <Separator />
              <div className="flex justify-between font-bold">
                <span>Total</span>
                <span>{formatRupiah(selected.total_amount)}</span>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  )
}
