"use client"

import { useEffect, useState } from "react"
import { toast } from "sonner"
import {
    Search,
    MoreVertical,
    FileText,
    Bike,
    Star,
    Plus,
    Users,
    Wifi,
    Clock,
    ShieldCheck,
    ActivitySquare,
    ChevronLeft,
    ChevronRight,
} from "lucide-react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
    Card,
    CardContent,
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
import { api } from "@/lib/api"
import { DriverFormModal } from "@/components/driver-form-modal"

interface Driver {
    id: string
    name: string
    email: string
    phone: string
    nik: string
    birth_place: string
    birth_date: string
    bank?: {
        id: number
        name: string
        code: string
    }
    bank_account_number?: string
    photo_profile_url?: string
    photo_ktp_url?: string
    photo_face_url?: string
    verification_status: "pending" | "accepted" | "rejected"
    is_active: boolean
    created_at: string
    updated_at: string
}

interface PaginationMeta {
    page: number
    limit: number
    total: number
    total_pages: number
}

const StatCard = ({
    title,
    value,
    subtitle,
    icon: Icon,
}: {
    title: string
    value: React.ReactNode
    subtitle: string
    icon?: any
}) => (
    <Card className="overflow-hidden">
        <div className="flex items-stretch h-full">
            <div className="w-1.5 bg-cakli-orange rounded-full my-3 ml-3 shrink-0" />
            <div className="flex-1">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-1 pt-3 px-3">
                    <CardTitle className="text-xs font-medium">{title}</CardTitle>
                    {Icon && <Icon className="h-3.5 w-3.5 text-muted-foreground" />}
                </CardHeader>
                <CardContent className="px-3 pb-3">
                    <div className="text-xl font-bold">{value}</div>
                    <p className="text-[10px] text-muted-foreground">{subtitle}</p>
                </CardContent>
            </div>
        </div>
    </Card>
)

export default function DriversPage() {
    const [drivers, setDrivers] = useState<Driver[]>([])
    const [loading, setLoading] = useState(true)
    const [searchQuery, setSearchQuery] = useState("")
    const [statusFilter, setStatusFilter] = useState("all")
    const [page, setPage] = useState(1)
    const [meta, setMeta] = useState<PaginationMeta>({
        page: 1,
        limit: 10,
        total: 0,
        total_pages: 0,
    })
    const [formModalOpen, setFormModalOpen] = useState(false)
    const [selectedDriver, setSelectedDriver] = useState<Driver | undefined>()

    useEffect(() => {
        fetchDrivers()
    }, [page, statusFilter, searchQuery])

    const fetchDrivers = async () => {
        try {
            setLoading(true)
            const params: any = {
                page,
                limit: 10,
            }
            
            if (searchQuery) params.search = searchQuery
            if (statusFilter !== "all") params.verification_status = statusFilter

            const response = await api.get("/admin/drivers", { params })
            setDrivers(response.data.data)
            setMeta(response.data.meta)
        } catch (error: any) {
            toast.error(error.response?.data?.message || "Gagal memuat data driver")
        } finally {
            setLoading(false)
        }
    }

    const getStatusBadge = (status: string) => {
        switch (status) {
            case "accepted":
                return <Badge variant="success" className="px-2 py-0.5 text-xs">Verified</Badge>
            case "pending":
                return <Badge variant="warning" className="px-2 py-0.5 text-xs">Pending</Badge>
            case "rejected":
                return <Badge variant="danger" className="px-2 py-0.5 text-xs">Rejected</Badge>
            default:
                return <Badge variant="outline" className="px-2 py-0.5 text-xs">{status}</Badge>
        }
    }

    const getInitials = (name: string) => {
        return name.split(" ").map(n => n[0]).join("").toUpperCase().substring(0, 2)
    }

    const handleSearch = (value: string) => {
        setSearchQuery(value)
        setPage(1)
    }

    const handleStatusFilter = (value: string) => {
        setStatusFilter(value)
        setPage(1)
    }

    if (loading && drivers.length === 0) {
        return (
            <div className="min-h-screen p-4 md:p-6 bg-white">
                <div className="max-w-[1400px] mx-auto space-y-4">
                    <div className="flex items-center justify-between">
                        <div>
                            <h1 className="text-2xl font-bold text-gray-900">Manajemen Driver</h1>
                            <p className="text-sm text-gray-500 mt-1">Loading...</p>
                        </div>
                    </div>
                    <div className="grid grid-cols-2 md:grid-cols-5 gap-3">
                        {[1, 2, 3, 4, 5].map((i) => (
                            <Card key={i} className="overflow-hidden">
                                <div className="flex items-stretch h-full">
                                    <div className="w-1.5 bg-gray-200 rounded-full my-3 ml-3 shrink-0 animate-pulse" />
                                    <div className="flex-1 p-3">
                                        <div className="h-4 w-24 bg-gray-200 rounded animate-pulse mb-2" />
                                        <div className="h-6 w-16 bg-gray-200 rounded animate-pulse mb-1" />
                                        <div className="h-3 w-32 bg-gray-200 rounded animate-pulse" />
                                    </div>
                                </div>
                            </Card>
                        ))}
                    </div>
                </div>
            </div>
        )
    }

    const stats = {
        total: meta.total,
        verified: drivers.filter(d => d.verification_status === "accepted").length,
        pending: drivers.filter(d => d.verification_status === "pending").length,
        rejected: drivers.filter(d => d.verification_status === "rejected").length,
        active: drivers.filter(d => d.is_active).length,
    }

    return (
        <div className="min-h-screen p-4 md:p-6 bg-white">
            <div className="max-w-[1400px] mx-auto space-y-4">
                {/* Header */}
                <div className="flex flex-col sm:flex-row sm:items-start justify-between gap-4">
                    <div>
                        <h1 className="text-2xl font-bold text-gray-900">Manajemen Driver</h1>
                        <p className="text-sm text-gray-500 mt-1">
                            Kelola driver, verifikasi dokumen, dan monitor performa armada.
                        </p>
                    </div>
                    <Button 
                        className="bg-[#E04D04] hover:bg-[#c94504] text-white"
                        onClick={() => {
                            setSelectedDriver(undefined)
                            setFormModalOpen(true)
                        }}
                    >
                        <Plus className="mr-2 h-4 w-4" />
                        Tambah Driver Baru
                    </Button>
                </div>

                {/* Stats Cards */}
                <div className="grid grid-cols-2 md:grid-cols-5 gap-3">
                    <StatCard
                        title="Total Driver"
                        value={stats.total}
                        subtitle="Registered drivers"
                        icon={Users}
                    />
                    <StatCard
                        title="Verified"
                        value={stats.verified}
                        subtitle={`${stats.active} active`}
                        icon={ShieldCheck}
                    />
                    <StatCard
                        title="Pending"
                        value={stats.pending}
                        subtitle="Menunggu verifikasi"
                        icon={Clock}
                    />
                    <StatCard
                        title="Rejected"
                        value={stats.rejected}
                        subtitle="Ditolak"
                        icon={ActivitySquare}
                    />
                    <StatCard
                        title="Online"
                        value="-"
                        subtitle="Real-time status"
                        icon={Wifi}
                    />
                </div>

                {/* Filter */}
                <div className="flex flex-col sm:flex-row items-start sm:items-center gap-3">
                    <div className="relative w-full sm:w-auto sm:flex-1 max-w-md">
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-600" />
                        <Input
                            type="search"
                            placeholder="Cari nama, email, NIK, atau telepon..."
                            className="pl-10 bg-white w-full h-11"
                            value={searchQuery}
                            onChange={(e) => handleSearch(e.target.value)}
                        />
                    </div>
                    <Select value={statusFilter} onValueChange={handleStatusFilter}>
                        <SelectTrigger className="w-full sm:w-40 bg-white h-11">
                            <SelectValue placeholder="Status" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="all">Semua Status</SelectItem>
                            <SelectItem value="accepted">Verified</SelectItem>
                            <SelectItem value="pending">Pending</SelectItem>
                            <SelectItem value="rejected">Rejected</SelectItem>
                        </SelectContent>
                    </Select>
                </div>

                {/* Table */}
                <div className="bg-white rounded-xl border border-slate-200 overflow-x-auto">
                    <Table>
                        <TableHeader>
                            <TableRow className="border-b border-gray-100">
                                <TableHead className="py-5 px-4 text-[11px] font-black text-slate-500 uppercase">Driver</TableHead>
                                <TableHead className="py-5 px-4 text-[11px] font-black text-slate-500 uppercase">NIK</TableHead>
                                <TableHead className="py-5 px-4 text-[11px] font-black text-slate-500 uppercase">Kontak</TableHead>
                                <TableHead className="py-5 px-4 text-[11px] font-black text-slate-500 uppercase">Bank</TableHead>
                                <TableHead className="py-5 px-4 text-[11px] font-black text-slate-500 uppercase">Status</TableHead>
                                <TableHead className="py-5 px-4 text-[11px] font-black text-slate-500 uppercase text-right">Aksi</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {drivers.map((driver) => (
                                <TableRow
                                    key={driver.id}
                                    className="border-b border-gray-100 hover:bg-slate-50/50"
                                >
                                    <TableCell className="py-4 px-4">
                                        <div className="flex items-center gap-3">
                                            <div className="w-10 h-10 rounded-xl bg-slate-100 flex items-center justify-center text-slate-700 font-black text-xs flex-shrink-0 border border-slate-200">
                                                {getInitials(driver.name)}
                                            </div>
                                            <div className="min-w-0">
                                                <div className="font-semibold text-slate-900 text-sm truncate">{driver.name}</div>
                                                <div className="text-[10px] text-slate-400 font-mono">{driver.email}</div>
                                            </div>
                                        </div>
                                    </TableCell>
                                    <TableCell className="py-4 px-4">
                                        <span className="text-sm font-mono text-slate-600">{driver.nik}</span>
                                    </TableCell>
                                    <TableCell className="py-4 px-4">
                                        <span className="text-sm text-slate-600">{driver.phone}</span>
                                    </TableCell>
                                    <TableCell className="py-4 px-4">
                                        {driver.bank ? (
                                            <div className="text-sm">
                                                <div className="font-semibold text-slate-900">{driver.bank.name}</div>
                                                <div className="text-xs text-slate-500">{driver.bank_account_number}</div>
                                            </div>
                                        ) : (
                                            <span className="text-slate-300">-</span>
                                        )}
                                    </TableCell>
                                    <TableCell className="py-4 px-4">{getStatusBadge(driver.verification_status)}</TableCell>
                                    <TableCell className="py-4 px-4 text-right">
                                        <div className="flex justify-end gap-2">
                                            <Button
                                                variant="ghost"
                                                size="sm"
                                                className="h-9 w-9 p-0 rounded-lg hover:bg-slate-100"
                                            >
                                                <FileText className="h-4 w-4 text-slate-600" />
                                            </Button>
                                            <DropdownMenu>
                                                <DropdownMenuTrigger asChild>
                                                    <Button variant="ghost" size="sm" className="h-9 w-9 p-0 rounded-lg hover:bg-slate-100">
                                                        <MoreVertical className="h-4 w-4 text-slate-600" />
                                                    </Button>
                                                </DropdownMenuTrigger>
                                                <DropdownMenuContent align="end" className="w-56">
                                                    <DropdownMenuLabel>Aksi</DropdownMenuLabel>
                                                    <DropdownMenuSeparator />
                                                    <DropdownMenuItem onClick={() => {
                                                        setSelectedDriver(driver)
                                                        setFormModalOpen(true)
                                                    }}>
                                                        Edit Driver
                                                    </DropdownMenuItem>
                                                    <DropdownMenuItem>Lihat Detail</DropdownMenuItem>
                                                    <DropdownMenuItem className="text-red-600">Hapus</DropdownMenuItem>
                                                </DropdownMenuContent>
                                            </DropdownMenu>
                                        </div>
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </div>

                {/* Pagination */}
                <div className="flex items-center justify-between">
                    <div className="text-sm text-slate-600">
                        Menampilkan {((page - 1) * meta.limit) + 1} - {Math.min(page * meta.limit, meta.total)} dari {meta.total} driver
                    </div>
                    <div className="flex items-center gap-2">
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => setPage(p => Math.max(1, p - 1))}
                            disabled={page === 1}
                        >
                            <ChevronLeft className="h-4 w-4" />
                        </Button>
                        <div className="text-sm font-medium">
                            Halaman {page} dari {meta.total_pages}
                        </div>
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => setPage(p => Math.min(meta.total_pages, p + 1))}
                            disabled={page === meta.total_pages}
                        >
                            <ChevronRight className="h-4 w-4" />
                        </Button>
                    </div>
                </div>

                {/* Form Modal */}
                <DriverFormModal
                    open={formModalOpen}
                    onOpenChange={setFormModalOpen}
                    driver={selectedDriver}
                    onSuccess={fetchDrivers}
                />
            </div>
        </div>
    )
}
