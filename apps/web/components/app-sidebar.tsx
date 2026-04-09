"use client"

import * as React from "react"
import { usePathname, useRouter } from "next/navigation"
import {
  LayoutDashboard,
  ShoppingCart,
  Users,
  Car,
  MessageSquareWarning,
  TrendingUp,
  Map,
  ShieldCheck,
  ClipboardList,
  TrendingDown,
  History,
  FileText,
  Home,
  Zap,
  EllipsisVertical,
  UserCheck,
  UserIcon,
  KeyIcon,
  LogOutIcon,
} from "lucide-react"

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { useAuth } from "@/lib/auth"
import { toast } from "sonner"

// ── Navigation data ──────────────────────────────────────────────────────────
const adminNavigation = {
  "reporting-admin": {
    items: [
      { title: "Dashboard", url: "/reporting-admin", icon: Home, exact: true },
      { title: "Analitik", url: "/reporting-admin/analytics", icon: TrendingUp },
      { title: "Pendapatan", url: "/reporting-admin/reports/revenue", icon: Zap },
      { title: "Riwayat Pesanan", url: "/reporting-admin/history", icon: History },
      { title: "Performa Pengemudi", url: "/reporting-admin/reports/driver-performance", icon: FileText },
      { title: "Laporan Pembatalan", url: "/reporting-admin/reports/cancellation", icon: TrendingDown },
    ],
  },
  "master-admin": {
    items: [
      { title: "Global Dashboard", url: "/master-admin", icon: LayoutDashboard, exact: true },
      { title: "Kelola Driver", url: "/master-admin/drivers", icon: Car },
      { title: "Kelola User", url: "/master-admin/users", icon: Users },
      { title: "Tariff Management", url: "/master-admin/tariffs", icon: TrendingUp },
      { title: "Area & Zone", url: "/master-admin/areas", icon: Map },
      { title: "Admin Roles", url: "/master-admin/roles", icon: ShieldCheck },
      { title: "Audit Log", url: "/master-admin/audit", icon: ClipboardList },
      { title: "Kebijakan Mitra", url: "/master-admin/partners", icon: UserCheck },
    ],
  },
  "operation-admin": {
    items: [
      { title: "Dashboard", url: "/operation-admin", icon: LayoutDashboard, exact: true },
      { title: "Live Map", url: "/operation-admin/map", icon: Map },
      { title: "Order Management", url: "/operation-admin/orders", icon: ShoppingCart },
      { title: "Driver Management", url: "/operation-admin/drivers", icon: Car },
      { title: "User Management", url: "/operation-admin/users", icon: Users },
      { title: "Activity Monitoring", url: "/operation-admin/activity", icon: Zap },
      { title: "Complaints & Disputes", url: "/operation-admin/complaints", icon: MessageSquareWarning },
    ],
  },
}

// ── AppSidebar ───────────────────────────────────────────────────────────────
export function AppSidebar() {
  const pathname = usePathname()
  const router = useRouter()
  const { admin, logout } = useAuth()

  const section = React.useMemo(() => {
    if (pathname.startsWith("/master-admin")) return "master-admin"
    if (pathname.startsWith("/operation-admin")) return "operation-admin"
    return "reporting-admin"
  }, [pathname])

  const { items } = adminNavigation[section]

  const handleLogout = async () => {
    try {
      await logout()
      toast.success("Berhasil logout")
      router.push("/")
    } catch (error) {
      toast.error("Gagal logout")
    }
  }

  const handleEditProfile = () => {
    toast.info("Fitur ubah profil akan segera hadir")
    // TODO: Implement edit profile modal
  }

  const handleChangePassword = () => {
    toast.info("Fitur ubah password akan segera hadir")
    // TODO: Implement change password modal
  }

  return (
    <aside
      style={{ width: "220px", minWidth: "220px", maxWidth: "220px" }}
      className="flex flex-col h-screen shrink-0"
    >
      {/* ── Logo / Header — uses public/cakli-logo.svg ── */}
      <div className="flex items-center px-4 pt-6 pb-6">
        <img
          src="/cakli-logo.svg"
          alt="CakLi"
          className="h-16 w-auto object-contain"
        />
      </div>

      {/* ── Navigation ── */}
      <nav className="flex flex-col gap-2 flex-1 overflow-y-auto">
        {items.map((item) => {
          const isActive = item.exact
            ? pathname === item.url
            : pathname.startsWith(item.url)

          if (isActive) {
            // Active: white pill with ONLY LEFT radius, right edge is square
            // This makes it visually "extend into" the white content rectangle
            return (
              <a
                key={item.title}
                href={item.url}
                className="flex items-center gap-3.5 ml-4 pl-4 pr-5 py-3 bg-[#FDFAF8] text-[#E65100] rounded-l-full font-semibold text-sm whitespace-nowrap transition-all"
              >
                <item.icon className="shrink-0 text-[#E65100]" style={{ width: 18, height: 18 }} strokeWidth={2.5} />
                <span>{item.title}</span>
              </a>
            )
          }

          return (
            <a
              key={item.title}
              href={item.url}
              className="flex items-center gap-3.5 mx-4 px-4 py-3 text-white/85 rounded-full font-medium text-sm whitespace-nowrap hover:bg-white/15 hover:text-white transition-all"
            >
              <item.icon className="shrink-0 text-white/80" style={{ width: 18, height: 18 }} strokeWidth={2} />
              <span>{item.title}</span>
            </a>
          )
        })}
      </nav>

      {/* ── User Footer ── */}
      <div className="mx-4 mb-5 mt-3 pt-3 border-t border-white/20">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <div className="flex items-center gap-2.5 px-2 py-2 rounded-2xl hover:bg-white/10 transition-colors cursor-pointer group">
              <Avatar className="h-8 w-8 rounded-xl shrink-0">
                <AvatarFallback className="rounded-xl bg-white/25 text-white text-xs font-bold">
                  {admin?.name?.substring(0, 2).toUpperCase() || "AD"}
                </AvatarFallback>
              </Avatar>
              <div className="flex flex-col flex-1 min-w-0 leading-tight">
                <span className="text-white text-xs font-semibold truncate">{admin?.name || "Admin"}</span>
                <span className="text-white/55 text-[10px] truncate">{admin?.email || "admin@cakli.id"}</span>
              </div>
              <EllipsisVertical className="w-4 h-4 text-white/50 shrink-0 opacity-0 group-hover:opacity-100 transition-opacity" />
            </div>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-56 rounded-lg"
            side="right"
            align="end"
            sideOffset={8}
          >
            <DropdownMenuLabel className="p-0 font-normal">
              <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                <Avatar className="h-8 w-8 rounded-lg">
                  <AvatarFallback className="rounded-lg bg-slate-100 text-slate-700 font-bold text-xs">
                    {admin?.name?.substring(0, 2).toUpperCase() || "AD"}
                  </AvatarFallback>
                </Avatar>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-medium">{admin?.name || "Admin"}</span>
                  <span className="text-muted-foreground truncate text-xs">
                    {admin?.email || "admin@cakli.id"}
                  </span>
                </div>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DropdownMenuItem onClick={handleEditProfile}>
                <UserIcon className="mr-2 h-4 w-4" />
                Ubah Profil
              </DropdownMenuItem>
              <DropdownMenuItem onClick={handleChangePassword}>
                <KeyIcon className="mr-2 h-4 w-4" />
                Ubah Password
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={handleLogout} className="text-red-600 focus:text-red-600 focus:bg-red-50">
              <LogOutIcon className="mr-2 h-4 w-4" />
              Keluar
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </aside>
  )
}
