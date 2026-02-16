import { BrowserRouter, Routes, Route } from "react-router-dom"
import { Toaster } from "@/components/ui/sonner"
import Layout from "@/components/Layout"
import Kasir from "@/pages/Kasir"
import Products from "@/pages/Products"
import Categories from "@/pages/Categories"
import Transactions from "@/pages/Transactions"
import Report from "@/pages/Report"

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<Layout />}>
          <Route path="/" element={<Kasir />} />
          <Route path="/products" element={<Products />} />
          <Route path="/categories" element={<Categories />} />
          <Route path="/transactions" element={<Transactions />} />
          <Route path="/report" element={<Report />} />
        </Route>
      </Routes>
      <Toaster richColors position="top-right" />
    </BrowserRouter>
  )
}
