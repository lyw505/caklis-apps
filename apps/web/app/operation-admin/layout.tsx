"use client"

import React, { useEffect } from "react"
import { useRouter } from "next/navigation"
import { AppSidebar } from "@/components/app-sidebar"
import { useAuth } from "@/lib/auth"

export default function OperationAdminLayout({
    children,
}: {
    children: React.ReactNode
}) {
    const { isAuthenticated, isLoading } = useAuth()
    const router = useRouter()

    useEffect(() => {
        if (!isLoading && !isAuthenticated) {
            router.push("/")
        }
    }, [isAuthenticated, isLoading, router])

    if (isLoading) {
        return (
            <div className="flex h-screen items-center justify-center">
                <div className="text-center">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-[#E65100] mx-auto"></div>
                    <p className="mt-4 text-gray-600">Loading...</p>
                </div>
            </div>
        )
    }

    if (!isAuthenticated) {
        return null
    }

    return (
        /* Full-screen orange background */
        <div className="flex h-screen overflow-hidden bg-[#E65100]">
            {/* Sidebar renders directly on orange bg */}
            <AppSidebar />

            {/* White content rectangle — radius only on top-left & bottom-left */}
            <div className="flex flex-col flex-1 overflow-hidden bg-white rounded-tl-[28px] rounded-bl-[28px]">

                <main className="flex-1 overflow-auto p-6 bg-[#FDFAF8]">
                    {children}
                </main>
            </div>
        </div>
    )
}