import { BrowserRouter, Routes, Route } from "react-router-dom"
import { Toaster } from "@/components/ui/sonner"
import Layout from "@/components/Layout"
import ProtectedRoute from "@/components/ProtectedRoute"
import Login from "@/pages/Login"
import Kasir from "@/pages/Kasir"
import Products from "@/pages/Products"
import Categories from "@/pages/Categories"
import Transactions from "@/pages/Transactions"
import Report from "@/pages/Report"

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route element={<ProtectedRoute><Layout /></ProtectedRoute>}>
          <Route path="/" element={<Kasir />} />
          <Route path="/products" element={
            <ProtectedRoute allowedRoles={["admin"]}><Products /></ProtectedRoute>
          } />
          <Route path="/categories" element={
            <ProtectedRoute allowedRoles={["admin"]}><Categories /></ProtectedRoute>
          } />
          <Route path="/transactions" element={<Transactions />} />
          <Route path="/report" element={<Report />} />
        </Route>
      </Routes>
      <Toaster richColors position="top-right" />
    </BrowserRouter>
  )
}
