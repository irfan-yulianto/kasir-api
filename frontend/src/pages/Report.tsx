import { useEffect, useState } from "react"
import { api, type SalesSummary } from "@/lib/api"
import { formatRupiah } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { toast } from "sonner"
import { Loader2, BarChart3, DollarSign, ShoppingCart, TrendingUp } from "lucide-react"

export default function Report() {
  const [report, setReport] = useState<SalesSummary | null>(null)
  const [loading, setLoading] = useState(true)
  const [startDate, setStartDate] = useState("")
  const [endDate, setEndDate] = useState("")

  const fetchReport = async () => {
    setLoading(true)
    try {
      const data = await api.getReport(startDate || undefined, endDate || undefined)
      setReport(data)
    } catch {
      toast.error("Gagal memuat laporan")
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchReport() }, [])

  const handleFilter = () => { fetchReport() }

  return (
    <div className="flex-1 overflow-y-auto p-6">
      <h1 className="text-2xl font-bold mb-6">Laporan Penjualan</h1>

      <div className="flex gap-3 mb-6 items-end">
        <div>
          <label className="text-sm font-medium">Dari Tanggal</label>
          <Input type="date" value={startDate} onChange={(e) => setStartDate(e.target.value)} />
        </div>
        <div>
          <label className="text-sm font-medium">Sampai Tanggal</label>
          <Input type="date" value={endDate} onChange={(e) => setEndDate(e.target.value)} />
        </div>
        <Button onClick={handleFilter}>Filter</Button>
        {(startDate || endDate) && (
          <Button variant="ghost" onClick={() => { setStartDate(""); setEndDate(""); setTimeout(fetchReport, 0) }}>Reset</Button>
        )}
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-20">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : !report ? (
        <div className="flex flex-col items-center justify-center py-20 text-muted-foreground">
          <BarChart3 className="h-12 w-12 mb-4" />
          <p>Tidak ada data laporan</p>
        </div>
      ) : (
        <>
          {/* Summary Cards */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 mb-8">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">Total Pendapatan</CardTitle>
                <DollarSign className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-bold">{formatRupiah(report.total_revenue)}</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">Total Transaksi</CardTitle>
                <ShoppingCart className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-bold">{report.total_transactions}</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">Rata-rata / Transaksi</CardTitle>
                <TrendingUp className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-bold">
                  {report.total_transactions > 0
                    ? formatRupiah(Math.round(report.total_revenue / report.total_transactions))
                    : formatRupiah(0)}
                </p>
              </CardContent>
            </Card>
          </div>

          {/* Top Products Table */}
          <h2 className="text-lg font-semibold mb-3">Top Produk Terlaris</h2>
          {report.top_products && report.top_products.length > 0 ? (
            <div className="border rounded-lg overflow-hidden">
              <table className="w-full text-sm">
                <thead className="bg-muted/50">
                  <tr>
                    <th className="text-left p-3 font-medium">Peringkat</th>
                    <th className="text-left p-3 font-medium">Produk</th>
                    <th className="text-right p-3 font-medium">Total Terjual</th>
                  </tr>
                </thead>
                <tbody>
                  {report.top_products.map((p, i) => (
                    <tr key={p.product_id} className="border-t">
                      <td className="p-3"><Badge variant="outline">#{i + 1}</Badge></td>
                      <td className="p-3 font-medium">{p.product_name}</td>
                      <td className="p-3 text-right">{p.total_sold} item</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <p className="text-muted-foreground text-sm">Belum ada data produk terjual</p>
          )}
        </>
      )}
    </div>
  )
}
