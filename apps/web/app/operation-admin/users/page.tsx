"use client"

import { useEffect, useState } from "react"
import { toast } from "sonner"
import {
    Search,
    MoreVertical,
    FileText,
    Plus,
    Users as UsersIcon,
    UserCheck,
    ChevronLeft,
    ChevronRight,
} from "lucide-react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
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
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { UserFormModal } from "@/components/user-form-modal"

interface User {
    id: string
    name: string
    email: string
    phone: string
    photo_profile_url?: string
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

export default function UsersPage() {
    const [users, setUsers] = useState<User[]>([])
    const [loading, setLoading] = useState(true)
    const [searchQuery, setSearchQuery] = useState("")
    const [page, setPage] = useState(1)
    const [meta, setMeta] = useState<PaginationMeta>({
        page: 1,
        limit: 10,
        total: 0,
        total_pages: 0,
    })
    const [formModalOpen, setFormModalOpen] = useState(false)
    const [selectedUser, setSelectedUser] = useState<User | undefined>()

    useEffect(() => {
        fetchUsers()
    }, [page, searchQuery])

    const fetchUsers = async () => {
        try {
            setLoading(true)
            const params: any = {
                page,
                limit: 10,
            }
            
            if (searchQuery) params.search = searchQuery

            const response = await api.get("/admin/users", { params })
            setUsers(response.data.data)
            setMeta(response.data.meta)
        } catch (error: any) {
            toast.error(error.response?.data?.message || "Gagal memuat data user")
        } finally {
            setLoading(false)
        }
    }

    const getInitials = (name: string) => {
        return name.split(" ").map(n => n[0]).join("").toUpperCase().substring(0, 2)
    }

    const handleSearch = (value: string) => {
        setSearchQuery(value)
        setPage(1)
    }

    if (loading && users.length === 0) {
        return (
            <div className="min-h-screen p-4 md:p-6 bg-white">
                <div className="max-w-[1400px] mx-auto space-y-4">
                    <div className="flex items-center justify-between">
                        <div>
                            <h1 className="text-2xl font-bold text-gray-900">Manajemen User</h1>
                            <p className="text-sm text-gray-500 mt-1">Loading...</p>
                        </div>
                    </div>
                    <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                        {[1, 2, 3].map((i) => (
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
        active: users.filter(u => u.is_active).length,
        inactive: users.filter(u => !u.is_active).length,
    }

    return (
        <div className="min-h-screen p-4 md:p-6 bg-white">
            <div className="max-w-[1400px] mx-auto space-y-4">
                {/* Header */}
                <div className="flex flex-col sm:flex-row sm:items-start justify-between gap-4">
                    <div>
                        <h1 className="text-2xl font-bold text-gray-900">Manajemen User</h1>
                        <p className="text-sm text-gray-500 mt-1">
                            Kelola pengguna aplikasi dan monitor aktivitas.
                        </p>
                    </div>
                    <Button 
                        className="bg-[#E04D04] hover:bg-[#c94504] text-white"
                        onClick={() => {
                            setSelectedUser(undefined)
                            setFormModalOpen(true)
                        }}
                    >
                        <Plus className="mr-2 h-4 w-4" />
                        Tambah User Baru
                    </Button>
                </div>

                {/* Stats Cards */}
                <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                    <StatCard
                        title="Total User"
                        value={stats.total}
                        subtitle="Registered users"
                        icon={UsersIcon}
                    />
                    <StatCard
                        title="Active"
                        value={stats.active}
                        subtitle="Currently active"
                        icon={UserCheck}
                    />
                    <StatCard
                        title="Inactive"
                        value={stats.inactive}
                        subtitle="Deactivated users"
                        icon={UsersIcon}
                    />
                </div>

                {/* Filter */}
                <div className="flex flex-col sm:flex-row items-start sm:items-center gap-3">
                    <div className="relative w-full sm:w-auto sm:flex-1 max-w-md">
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-600" />
                        <Input
                            type="search"
                            placeholder="Cari nama, email, atau telepon..."
                            className="pl-10 bg-white w-full h-11"
                            value={searchQuery}
                            onChange={(e) => handleSearch(e.target.value)}
                        />
                    </div>
                </div>

                {/* Table */}
                <div className="bg-white rounded-xl border border-slate-200 overflow-x-auto">
                    <Table>
                        <TableHeader>
                            <TableRow className="border-b border-gray-100">
                                <TableHead className="py-5 px-4 text-[11px] font-black text-slate-500 uppercase">User</TableHead>
                                <TableHead className="py-5 px-4 text-[11px] font-black text-slate-500 uppercase">Email</TableHead>
                                <TableHead className="py-5 px-4 text-[11px] font-black text-slate-500 uppercase">Telepon</TableHead>
                                <TableHead className="py-5 px-4 text-[11px] font-black text-slate-500 uppercase">Status</TableHead>
                                <TableHead className="py-5 px-4 text-[11px] font-black text-slate-500 uppercase">Bergabung</TableHead>
                                <TableHead className="py-5 px-4 text-[11px] font-black text-slate-500 uppercase text-right">Aksi</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {users.map((user) => (
                                <TableRow
                                    key={user.id}
                                    className="border-b border-gray-100 hover:bg-slate-50/50"
                                >
                                    <TableCell className="py-4 px-4">
                                        <div className="flex items-center gap-3">
                                            <Avatar className="h-10 w-10">
                                                <AvatarImage src={user.photo_profile_url} alt={user.name} />
                                                <AvatarFallback className="bg-slate-100 text-slate-700 font-bold text-xs">
                                                    {getInitials(user.name)}
                                                </AvatarFallback>
                                            </Avatar>
                                            <div className="min-w-0">
                                                <div className="font-semibold text-slate-900 text-sm truncate">{user.name}</div>
                                                <div className="text-[10px] text-slate-400 font-mono">{user.id}</div>
                                            </div>
                                        </div>
                                    </TableCell>
                                    <TableCell className="py-4 px-4">
                                        <span className="text-sm text-slate-600">{user.email}</span>
                                    </TableCell>
                                    <TableCell className="py-4 px-4">
                                        <span className="text-sm text-slate-600">{user.phone}</span>
                                    </TableCell>
                                    <TableCell className="py-4 px-4">
                                        {user.is_active ? (
                                            <Badge variant="success" className="px-2 py-0.5 text-xs">Active</Badge>
                                        ) : (
                                            <Badge variant="neutral" className="px-2 py-0.5 text-xs">Inactive</Badge>
                                        )}
                                    </TableCell>
                                    <TableCell className="py-4 px-4">
                                        <span className="text-sm text-slate-600">
                                            {new Date(user.created_at).toLocaleDateString('id-ID', {
                                                year: 'numeric',
                                                month: 'short',
                                                day: 'numeric'
                                            })}
                                        </span>
                                    </TableCell>
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
                                                        setSelectedUser(user)
                                                        setFormModalOpen(true)
                                                    }}>
                                                        Edit User
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
                        Menampilkan {((page - 1) * meta.limit) + 1} - {Math.min(page * meta.limit, meta.total)} dari {meta.total} user
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
                <UserFormModal
                    open={formModalOpen}
                    onOpenChange={setFormModalOpen}
                    user={selectedUser}
                    onSuccess={fetchUsers}
                />
            </div>
        </div>
    )
}
