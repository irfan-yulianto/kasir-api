import { NavLink, Outlet } from "react-router-dom"
import { cn } from "@/lib/utils"
import {
  ShoppingCart, Package, Tag, Receipt, BarChart3,
} from "lucide-react"

const navItems = [
  { to: "/", label: "Kasir", icon: ShoppingCart },
  { to: "/products", label: "Produk", icon: Package, group: "Master Data" },
  { to: "/categories", label: "Kategori", icon: Tag },
  { to: "/transactions", label: "Riwayat Transaksi", icon: Receipt, group: "Transaksi" },
  { to: "/report", label: "Penjualan", icon: BarChart3, group: "Laporan" },
]

export default function Layout() {
  let lastGroup = ""

  return (
    <div className="flex h-screen bg-background">
      {/* Sidebar */}
      <aside className="w-56 border-r bg-card flex flex-col shrink-0">
        <div className="p-4 border-b">
          <h1 className="text-lg font-bold">Kasir App</h1>
        </div>
        <nav className="flex-1 p-3 space-y-1">
          {navItems.map((item) => {
            const showGroup = item.group && item.group !== lastGroup
            if (item.group) lastGroup = item.group
            return (
              <div key={item.to}>
                {showGroup && (
                  <p className="text-xs font-medium text-muted-foreground uppercase tracking-wider px-3 pt-4 pb-1">
                    {item.group}
                  </p>
                )}
                <NavLink
                  to={item.to}
                  end={item.to === "/"}
                  className={({ isActive }) =>
                    cn(
                      "flex items-center gap-3 rounded-md px-3 py-2 text-sm transition-colors",
                      isActive
                        ? "bg-primary text-primary-foreground"
                        : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
                    )
                  }
                >
                  <item.icon className="h-4 w-4" />
                  {item.label}
                </NavLink>
              </div>
            )
          })}
        </nav>
      </aside>

      {/* Main */}
      <div className="flex-1 flex flex-col overflow-hidden">
        <Outlet />
      </div>
    </div>
  )
}
