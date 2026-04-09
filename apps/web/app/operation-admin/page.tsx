"use client"

import { useEffect, useState } from "react"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { AlertCircle, Clock, MapPin, CheckCircle2, MessageSquareWarning, ArrowUpRight, Activity, Users, Car } from "lucide-react"
import { ChartAreaInteractive } from "@/components/chart-area-interactive"
import { api } from "@/lib/api"
import { toast } from "sonner"

interface DashboardStats {
  drivers: {
    total: number
    verified: number
    pending: number
    rejected: number
  }
  users: {
    total: number
    active: number
  }
}

export default function DashboardPage() {
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchStats()
  }, [])

  const fetchStats = async () => {
    try {
      setLoading(true)
      const response = await api.get("/admin/dashboard/stats")
      setStats(response.data.data)
    } catch (error: any) {
      toast.error(error.response?.data?.message || "Gagal memuat statistik")
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="flex flex-col gap-6 p-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Dashboard Operasional</h1>
            <p className="text-muted-foreground">Overview of current system status and activities.</p>
          </div>
        </div>
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
          {[1, 2, 3, 4].map((i) => (
            <Card key={i} className="overflow-hidden">
              <div className="flex items-stretch h-full">
                <div className="w-1.5 bg-gray-200 rounded-full my-3 ml-3 shrink-0 animate-pulse" />
                <div className="flex-1">
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-1 pt-3 px-3">
                    <div className="h-4 w-24 bg-gray-200 rounded animate-pulse" />
                  </CardHeader>
                  <CardContent className="px-3 pb-3">
                    <div className="h-6 w-16 bg-gray-200 rounded animate-pulse mb-2" />
                    <div className="h-3 w-32 bg-gray-200 rounded animate-pulse" />
                  </CardContent>
                </div>
              </div>
            </Card>
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-6 p-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Dashboard Operasional</h1>
          <p className="text-muted-foreground">Overview of current system status and activities.</p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline" className="text-green-600 border-green-600 bg-green-50 dark:bg-green-950/20 px-3 py-1">
            <Activity className="w-3 h-3 mr-2" />
            System Normal
          </Badge>
          <span className="text-sm text-muted-foreground">Last updated: Just now</span>
        </div>
      </div>

      {/* Quick Stats */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
        <Card className="overflow-hidden">
          <div className="flex items-stretch h-full">
            <div className="w-1.5 bg-cakli-orange rounded-full my-3 ml-3 shrink-0" />
            <div className="flex-1">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-1 pt-3 px-3">
                <CardTitle className="text-xs font-medium">Total Driver</CardTitle>
                <Car className="h-3.5 w-3.5 text-muted-foreground" />
              </CardHeader>
              <CardContent className="px-3 pb-3">
                <div className="text-xl font-bold">{stats?.drivers.total || 0}</div>
                <p className="text-[10px] text-muted-foreground">
                  {stats?.drivers.verified || 0} verified
                </p>
              </CardContent>
            </div>
          </div>
        </Card>

        <Card className="overflow-hidden">
          <div className="flex items-stretch h-full">
            <div className="w-1.5 bg-cakli-orange rounded-full my-3 ml-3 shrink-0" />
            <div className="flex-1">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-1 pt-3 px-3">
                <CardTitle className="text-xs font-medium">Driver Pending</CardTitle>
                <Clock className="h-3.5 w-3.5 text-muted-foreground" />
              </CardHeader>
              <CardContent className="px-3 pb-3">
                <div className="text-xl font-bold">{stats?.drivers.pending || 0}</div>
                <p className="text-[10px] text-muted-foreground">Menunggu verifikasi</p>
              </CardContent>
            </div>
          </div>
        </Card>

        <Card className="overflow-hidden">
          <div className="flex items-stretch h-full">
            <div className="w-1.5 bg-cakli-orange rounded-full my-3 ml-3 shrink-0" />
            <div className="flex-1">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-1 pt-3 px-3">
                <CardTitle className="text-xs font-medium">Total User</CardTitle>
                <Users className="h-3.5 w-3.5 text-muted-foreground" />
              </CardHeader>
              <CardContent className="px-3 pb-3">
                <div className="text-xl font-bold">{stats?.users.total || 0}</div>
                <p className="text-[10px] text-muted-foreground">
                  {stats?.users.active || 0} active
                </p>
              </CardContent>
            </div>
          </div>
        </Card>

        <Card className="overflow-hidden">
          <div className="flex items-stretch h-full">
            <div className="w-1.5 bg-cakli-orange rounded-full my-3 ml-3 shrink-0" />
            <div className="flex-1">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-1 pt-3 px-3">
                <CardTitle className="text-xs font-medium">Driver Rejected</CardTitle>
                <AlertCircle className="h-3.5 w-3.5 text-red-500" />
              </CardHeader>
              <CardContent className="px-3 pb-3">
                <div className="text-xl font-bold text-red-500">{stats?.drivers.rejected || 0}</div>
                <p className="text-[10px] text-muted-foreground">Ditolak</p>
              </CardContent>
            </div>
          </div>
        </Card>
      </div>

      <div className="grid gap-6 lg:grid-cols-7">
        {/* Main Content Area */}
        <div className="col-span-1 lg:col-span-4 flex flex-col gap-6">
          {/* Chart Section */}
          <ChartAreaInteractive />

          {/* Live Orders Table */}
          <Card>
            <CardHeader>
              <CardTitle>Live Orders</CardTitle>
              <CardDescription>Real-time status of ongoing trips.</CardDescription>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Order ID</TableHead>
                    <TableHead>Customer</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="text-right">Action</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {[
                    { id: "ORD-001", cust: "Rina S.", status: "Picking Up" },
                    { id: "ORD-002", cust: "Ahmad J.", status: "In Transit" },
                    { id: "ORD-003", cust: "Dewi P.", status: "In Transit" },
                    { id: "ORD-004", cust: "Kevin L.", status: "Assigning" },
                    { id: "ORD-005", cust: "Sarah M.", status: "Picking Up" },
                  ].map((order) => (
                    <TableRow key={order.id}>
                      <TableCell className="font-medium">{order.id}</TableCell>
                      <TableCell>{order.cust}</TableCell>
                      <TableCell>
                        <Badge variant={order.status === "Assigning" ? "secondary" : "default"}>
                          {order.status}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right">
                        <Button variant="ghost" size="sm">Details</Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </div>

        {/* Sidebar Content Area */}
        <div className="col-span-1 lg:grid-cols-3 flex flex-col gap-6">
          {/* Quick Actions */}
          <Card className="border-orange-200 bg-orange-50/30">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-bold uppercase tracking-wider text-orange-900">Quick Operations</CardTitle>
            </CardHeader>
            <CardContent className="grid grid-cols-2 gap-2">
              <Button variant="outline" className="h-auto py-3 px-2 flex flex-col gap-1 border-orange-200 hover:bg-orange-100 text-orange-800">
                <Activity className="h-4 w-4" />
                <span className="text-[10px] uppercase font-bold">Broadcast</span>
              </Button>
              <Button variant="outline" className="h-auto py-3 px-2 flex flex-col gap-1 border-orange-200 hover:bg-orange-100 text-orange-800">
                <MapPin className="h-4 w-4" />
                <span className="text-[10px] uppercase font-bold">Heatmap</span>
              </Button>
            </CardContent>
          </Card>

          {/* Alerts & Notifications */}
          <Card>
            <CardHeader>
              <CardTitle>System Alerts</CardTitle>
              <CardDescription>Critical updates needing attention.</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-start gap-4 rounded-lg border p-3 bg-red-50 dark:bg-red-950/20">
                <AlertCircle className="h-5 w-5 text-red-600 mt-0.5" />
                <div className="grid gap-1">
                  <p className="font-medium text-red-900 dark:text-red-200">Emergency Signal</p>
                  <p className="text-sm text-red-800 dark:text-red-300">
                    Driver "Agus T." triggered PANIC button.
                  </p>
                  <div className="flex gap-2">
                    <Button size="sm" variant="destructive" className="mt-1">Call Driver</Button>
                    <Button size="sm" variant="outline" className="mt-1 bg-white">Track Now</Button>
                  </div>
                </div>
              </div>
              <div className="flex items-start gap-4 rounded-lg border p-3 bg-orange-50 dark:bg-orange-950/20 border-orange-200">
                <MessageSquareWarning className="h-5 w-5 text-orange-600 mt-0.5" />
                <div className="grid gap-1">
                  <p className="font-medium text-orange-900 dark:text-orange-200">High Demand Alert</p>
                  <p className="text-sm text-orange-800 dark:text-orange-300">
                    Sudden surge in "Malang Kota" area.
                  </p>
                  <Button size="sm" variant="outline" className="w-fit mt-1 border-orange-200 hover:bg-orange-100 text-orange-700">Adjust Payout</Button>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Recent Activity */}
          <Card>
            <CardHeader>
              <CardTitle>Fleet Statistics</CardTitle>
              <CardDescription>Real-time vehicle status.</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="flex items-center justify-between text-sm">
                  <div className="flex items-center gap-2">
                    <div className="size-2 rounded-full bg-green-500" />
                    <span>Verified</span>
                  </div>
                  <span className="font-bold">{stats?.drivers.verified || 0}</span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <div className="flex items-center gap-2">
                    <div className="size-2 rounded-full bg-yellow-500" />
                    <span>Pending</span>
                  </div>
                  <span className="font-bold">{stats?.drivers.pending || 0}</span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <div className="flex items-center gap-2">
                    <div className="size-2 rounded-full bg-red-500" />
                    <span>Rejected</span>
                  </div>
                  <span className="font-bold">{stats?.drivers.rejected || 0}</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
