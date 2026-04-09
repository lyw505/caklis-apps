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
import { api } from "@/lib/api"

interface UserFormData {
    name: string
    email: string
    password?: string
    phone: string
    photo_profile_key?: string
}

interface UserFormModalProps {
    open: boolean
    onOpenChange: (open: boolean) => void
    user?: any
    onSuccess: () => void
}

export function UserFormModal({ open, onOpenChange, user, onSuccess }: UserFormModalProps) {
    const [loading, setLoading] = useState(false)
    const [formData, setFormData] = useState<UserFormData>({
        name: "",
        email: "",
        password: "",
        phone: "",
    })
    const [photoProfile, setPhotoProfile] = useState<File | null>(null)
    const [uploading, setUploading] = useState(false)

    useEffect(() => {
        if (open) {
            if (user) {
                setFormData({
                    name: user.name || "",
                    email: user.email || "",
                    phone: user.phone || "",
                    photo_profile_key: user.photo_profile_url,
                })
            } else {
                resetForm()
            }
        }
    }, [open, user])

    const resetForm = () => {
        setFormData({
            name: "",
            email: "",
            password: "",
            phone: "",
        })
        setPhotoProfile(null)
    }

    const uploadFile = async (file: File): Promise<string> => {
        try {
            // Get presigned upload URL
            const uploadResponse = await api.post("/upload/presigned-url", {
                filename: file.name,
                content_type: file.type,
                folder: "users/photo-profile",
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
        if (!formData.name || !formData.email || !formData.phone) {
            toast.error("Mohon lengkapi semua field yang wajib diisi")
            return
        }

        if (!user && !formData.password) {
            toast.error("Password wajib diisi untuk user baru")
            return
        }

        try {
            setLoading(true)
            setUploading(true)

            // Upload photo if selected
            if (photoProfile) {
                formData.photo_profile_key = await uploadFile(photoProfile)
            }

            setUploading(false)

            // Create or update user
            if (user) {
                await api.put(`/admin/users/${user.id}`, formData)
                toast.success("User berhasil diupdate")
            } else {
                await api.post("/admin/users", formData)
                toast.success("User berhasil ditambahkan")
            }

            onSuccess()
            onOpenChange(false)
        } catch (error: any) {
            toast.error(error.response?.data?.message || "Gagal menyimpan user")
        } finally {
            setLoading(false)
            setUploading(false)
        }
    }

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-md">
                <DialogHeader>
                    <DialogTitle>{user ? "Edit User" : "Tambah User Baru"}</DialogTitle>
                    <DialogDescription>
                        {user ? "Update informasi user" : "Lengkapi form untuk menambahkan user baru"}
                    </DialogDescription>
                </DialogHeader>

                <form onSubmit={handleSubmit} className="space-y-4">
                    {/* Photo Upload */}
                    <div className="space-y-2">
                        <Label>Foto Profile</Label>
                        <div className="border-2 border-dashed rounded-lg p-4 text-center">
                            {photoProfile ? (
                                <div className="relative inline-block">
                                    <img
                                        src={URL.createObjectURL(photoProfile)}
                                        alt="Preview"
                                        className="w-32 h-32 object-cover rounded-full mx-auto"
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
                                <label className="cursor-pointer block">
                                    <Upload className="h-12 w-12 mx-auto text-gray-400" />
                                    <span className="text-sm text-gray-500 mt-2 block">Upload Foto Profile</span>
                                    <span className="text-xs text-gray-400">JPG atau PNG, max 5MB</span>
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

                    {/* Basic Info */}
                    <div className="space-y-4">
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
                            <Label htmlFor="phone">No. Telepon *</Label>
                            <Input
                                id="phone"
                                value={formData.phone}
                                onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                                required
                            />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="password">Password {!user && "*"}</Label>
                            <Input
                                id="password"
                                type="password"
                                value={formData.password}
                                onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                                placeholder={user ? "Kosongkan jika tidak diubah" : ""}
                                required={!user}
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
