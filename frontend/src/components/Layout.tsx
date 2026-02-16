import { NavLink, Outlet, useNavigate } from "react-router-dom"
import { useAuth } from "@/contexts/AuthContext"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
  ShoppingCart, Package, Tag, Receipt, BarChart3, LogOut, User,
} from "lucide-react"

const navItems = [
  { to: "/", label: "Kasir", icon: ShoppingCart },
  { to: "/products", label: "Produk", icon: Package, group: "Master Data", roles: ["admin"] },
  { to: "/categories", label: "Kategori", icon: Tag, roles: ["admin"] },
  { to: "/transactions", label: "Riwayat Transaksi", icon: Receipt, group: "Transaksi" },
  { to: "/report", label: "Penjualan", icon: BarChart3, group: "Laporan" },
]

export default function Layout() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()
  let lastGroup = ""

  const handleLogout = () => {
    logout()
    navigate("/login")
  }

  const visibleItems = navItems.filter(
    (item) => !item.roles || (user && item.roles.includes(user.role))
  )

  return (
    <div className="flex h-screen bg-background">
      {/* Sidebar */}
      <aside className="w-56 border-r bg-card flex flex-col shrink-0">
        <div className="p-4 border-b">
          <h1 className="text-lg font-bold">Kasir App</h1>
        </div>
        <nav className="flex-1 p-3 space-y-1">
          {visibleItems.map((item) => {
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
        {/* User info + Logout */}
        <div className="border-t p-3 space-y-2">
          <div className="flex items-center gap-2 px-3 py-1">
            <User className="h-4 w-4 text-muted-foreground" />
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium truncate">{user?.username}</p>
              <p className="text-xs text-muted-foreground capitalize">{user?.role}</p>
            </div>
          </div>
          <Button variant="outline" size="sm" className="w-full" onClick={handleLogout}>
            <LogOut className="h-4 w-4 mr-2" />
            Logout
          </Button>
        </div>
      </aside>

      {/* Main */}
      <div className="flex-1 flex flex-col overflow-hidden">
        <Outlet />
      </div>
    </div>
  )
}
