"use client"

import { useEffect, useState } from "react"
import { toast } from "sonner"
import { Upload, X, Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import { api } from "@/lib/api"

interface Bank {
    id: number
    name: string
    code: string
}

interface DriverFormData {
    name: string
    email: string
    password?: string
    phone: string
    nik: string
    birth_place: string
    birth_date: string
    bank_id?: number
    bank_account_number?: string
    verification_status: string
    photo_profile_key?: string
    photo_ktp_key?: string
    photo_face_key?: string
}

interface DriverFormModalProps {
    open: boolean
    onOpenChange: (open: boolean) => void
    driver?: any
    onSuccess: () => void
}

export function DriverFormModal({ open, onOpenChange, driver, onSuccess }: DriverFormModalProps) {
    const [loading, setLoading] = useState(false)
    const [banks, setBanks] = useState<Bank[]>([])
    const [formData, setFormData] = useState<DriverFormData>({
        name: "",
        email: "",
        password: "",
        phone: "",
        nik: "",
        birth_place: "",
        birth_date: "",
        verification_status: "pending",
    })
    const [photoProfile, setPhotoProfile] = useState<File | null>(null)
    const [photoKTP, setPhotoKTP] = useState<File | null>(null)
    const [photoFace, setPhotoFace] = useState<File | null>(null)
    const [uploading, setUploading] = useState(false)

    useEffect(() => {
        if (open) {
            fetchBanks()
            if (driver) {
                setFormData({
                    name: driver.name || "",
                    email: driver.email || "",
                    phone: driver.phone || "",
                    nik: driver.nik || "",
                    birth_place: driver.birth_place || "",
                    birth_date: driver.birth_date || "",
                    bank_id: driver.bank?.id,
                    bank_account_number: driver.bank_account_number || "",
                    verification_status: driver.verification_status || "pending",
                    photo_profile_key: driver.photo_profile_url,
                    photo_ktp_key: driver.photo_ktp_url,
                    photo_face_key: driver.photo_face_url,
                })
            } else {
                resetForm()
            }
        }
    }, [open, driver])

    const fetchBanks = async () => {
        try {
            const response = await api.get("/admin/banks")
            setBanks(response.data.data)
        } catch (error: any) {
            toast.error("Gagal memuat data bank")
        }
    }

    const resetForm = () => {
        setFormData({
            name: "",
            email: "",
            password: "",
            phone: "",
            nik: "",
            birth_place: "",
            birth_date: "",
            verification_status: "pending",
        })
        setPhotoProfile(null)
        setPhotoKTP(null)
        setPhotoFace(null)
    }

    const uploadFile = async (file: File, folder: string): Promise<string> => {
        try {
            // Get presigned upload URL
            const uploadResponse = await api.post("/upload/presigned-url", {
                filename: file.name,
                content_type: file.type,
                folder: folder,
            })

            const { upload_url, object_key } = uploadResponse.data.data

            // Upload file to MinIO
            await fetch(upload_url, {
                method: "PUT",
                body: file,
                headers: {
                    "Content-Type": file.type,
                },
            })

            return object_key
        } catch (error: any) {
            throw new Error(error.response?.data?.message || "Gagal upload file")
        }
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        
        // Validation
        if (!formData.name || !formData.email || !formData.phone || !formData.nik) {
            toast.error("Mohon lengkapi semua field yang wajib diisi")
            return
        }

        if (!driver && !formData.password) {
            toast.error("Password wajib diisi untuk driver baru")
            return
        }

        if (formData.nik.length !== 16) {
            toast.error("NIK harus 16 digit")
            return
        }

        try {
            setLoading(true)
            setUploading(true)

            // Upload photos if selected
            if (photoProfile) {
                formData.photo_profile_key = await uploadFile(photoProfile, "drivers/photo-profile")
            }
            if (photoKTP) {
                formData.photo_ktp_key = await uploadFile(photoKTP, "drivers/photo-ktp")
            }
            if (photoFace) {
                formData.photo_face_key = await uploadFile(photoFace, "drivers/photo-face")
            }

            setUploading(false)

            // Create or update driver
            if (driver) {
                await api.put(`/admin/drivers/${driver.id}`, formData)
                toast.success("Driver berhasil diupdate")
            } else {
                await api.post("/admin/drivers", formData)
                toast.success("Driver berhasil ditambahkan")
            }

            onSuccess()
            onOpenChange(false)
        } catch (error: any) {
            toast.error(error.response?.data?.message || "Gagal menyimpan driver")
        } finally {
            setLoading(false)
            setUploading(false)
        }
    }

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle>{driver ? "Edit Driver" : "Tambah Driver Baru"}</DialogTitle>
                    <DialogDescription>
                        {driver ? "Update informasi driver" : "Lengkapi form untuk menambahkan driver baru"}
                    </DialogDescription>
                </DialogHeader>

                <form onSubmit={handleSubmit} className="space-y-4">
                    {/* Photo Uploads */}
                    <div className="grid grid-cols-3 gap-4">
                        <div className="space-y-2">
                            <Label>Foto Profile</Label>
                            <div className="border-2 border-dashed rounded-lg p-4 text-center">
                                {photoProfile ? (
                                    <div className="relative">
                                        <img
                                            src={URL.createObjectURL(photoProfile)}
                                            alt="Preview"
                                            className="w-full h-24 object-cover rounded"
                                        />
                                        <Button
                                            type="button"
                                            variant="destructive"
                                            size="sm"
                                            className="absolute top-0 right-0"
                                            onClick={() => setPhotoProfile(null)}
                                        >
                                            <X className="h-3 w-3" />
                                        </Button>
                                    </div>
                                ) : (
                                    <label className="cursor-pointer">
                                        <Upload className="h-8 w-8 mx-auto text-gray-400" />
                                        <span className="text-xs text-gray-500">Upload</span>
                                        <input
                                            type="file"
                                            accept="image/jpeg,image/png"
                                            className="hidden"
                                            onChange={(e) => setPhotoProfile(e.target.files?.[0] || null)}
                                        />
                                    </label>
                                )}
                            </div>
                        </div>

                        <div className="space-y-2">
                            <Label>Foto KTP</Label>
                            <div className="border-2 border-dashed rounded-lg p-4 text-center">
                                {photoKTP ? (
                                    <div className="relative">
                                        <img
                                            src={URL.createObjectURL(photoKTP)}
                                            alt="Preview"
                                            className="w-full h-24 object-cover rounded"
                                        />
                                        <Button
                                            type="button"
                                            variant="destructive"
                                            size="sm"
                                            className="absolute top-0 right-0"
                                            onClick={() => setPhotoKTP(null)}
                                        >
                                            <X className="h-3 w-3" />
                                        </Button>
                                    </div>
                                ) : (
                                    <label className="cursor-pointer">
                                        <Upload className="h-8 w-8 mx-auto text-gray-400" />
                                        <span className="text-xs text-gray-500">Upload</span>
                                        <input
                                            type="file"
                                            accept="image/jpeg,image/png"
                                            className="hidden"
                                            onChange={(e) => setPhotoKTP(e.target.files?.[0] || null)}
                                        />
                                    </label>
                                )}
                            </div>
                        </div>

                        <div className="space-y-2">
                            <Label>Foto Muka</Label>
                            <div className="border-2 border-dashed rounded-lg p-4 text-center">
                                {photoFace ? (
                                    <div className="relative">
                                        <img
                                            src={URL.createObjectURL(photoFace)}
                                            alt="Preview"
                                            className="w-full h-24 object-cover rounded"
                                        />
                                        <Button
                                            type="button"
                                            variant="destructive"
                                            size="sm"
                                            className="absolute top-0 right-0"
                                            onClick={() => setPhotoFace(null)}
                                        >
                                            <X className="h-3 w-3" />
                                        </Button>
                                    </div>
                                ) : (
                                    <label className="cursor-pointer">
                                        <Upload className="h-8 w-8 mx-auto text-gray-400" />
                                        <span className="text-xs text-gray-500">Upload</span>
                                        <input
                                            type="file"
                                            accept="image/jpeg,image/png"
                                            className="hidden"
                                            onChange={(e) => setPhotoFace(e.target.files?.[0] || null)}
                                        />
                                    </label>
                                )}
                            </div>
                        </div>
                    </div>

                    {/* Basic Info */}
                    <div className="grid grid-cols-2 gap-4">
                        <div className="space-y-2">
                            <Label htmlFor="name">Nama Lengkap *</Label>
                            <Input
                                id="name"
                                value={formData.name}
                                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                                required
                            />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="email">Email *</Label>
                            <Input
                                id="email"
                                type="email"
                                value={formData.email}
                                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                                required
                            />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="password">Password {!driver && "*"}</Label>
                            <Input
                                id="password"
                                type="password"
                                value={formData.password}
                                onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                                placeholder={driver ? "Kosongkan jika tidak diubah" : ""}
                                required={!driver}
                            />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="phone">No. Telepon *</Label>
                            <Input
                                id="phone"
                                value={formData.phone}
                                onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                                required
                            />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="nik">NIK (16 digit) *</Label>
                            <Input
                                id="nik"
                                value={formData.nik}
                                onChange={(e) => setFormData({ ...formData, nik: e.target.value })}
                                maxLength={16}
                                required
                            />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="birth_place">Tempat Lahir</Label>
                            <Input
                                id="birth_place"
                                value={formData.birth_place}
                                onChange={(e) => setFormData({ ...formData, birth_place: e.target.value })}
                            />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="birth_date">Tanggal Lahir</Label>
                            <Input
                                id="birth_date"
                                type="date"
                                value={formData.birth_date}
                                onChange={(e) => setFormData({ ...formData, birth_date: e.target.value })}
                            />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="verification_status">Status Verifikasi</Label>
                            <Select
                                value={formData.verification_status}
                                onValueChange={(value) => setFormData({ ...formData, verification_status: value })}
                            >
                                <SelectTrigger>
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="pending">Pending</SelectItem>
                                    <SelectItem value="accepted">Accepted</SelectItem>
                                    <SelectItem value="rejected">Rejected</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="bank_id">Bank</Label>
                            <Select
                                value={formData.bank_id?.toString()}
                                onValueChange={(value) => setFormData({ ...formData, bank_id: parseInt(value) })}
                            >
                                <SelectTrigger>
                                    <SelectValue placeholder="Pilih bank" />
                                </SelectTrigger>
                                <SelectContent>
                                    {banks.map((bank) => (
                                        <SelectItem key={bank.id} value={bank.id.toString()}>
                                            {bank.name}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="bank_account_number">No. Rekening</Label>
                            <Input
                                id="bank_account_number"
                                value={formData.bank_account_number}
                                onChange={(e) => setFormData({ ...formData, bank_account_number: e.target.value })}
                            />
                        </div>
                    </div>

                    <DialogFooter>
                        <Button type="button" variant="outline" onClick={() => onOpenChange(false)} disabled={loading}>
                            Batal
                        </Button>
                        <Button type="submit" disabled={loading} className="bg-[#E04D04] hover:bg-[#c94504]">
                            {loading ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    {uploading ? "Uploading..." : "Menyimpan..."}
                                </>
                            ) : (
                                "Simpan"
                            )}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    )
}
