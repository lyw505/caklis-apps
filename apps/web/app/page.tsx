"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Image from "next/image";
import { Eye, EyeOff } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useAuth } from "@/lib/auth";
import { toast } from "sonner";
import api from "@/lib/api";

export default function LoginPage() {
    const [showPassword, setShowPassword] = useState(false);
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [isLoading, setIsLoading] = useState(false);
    
    const { login, isAuthenticated, admin } = useAuth();
    const router = useRouter();

    // Redirect if already authenticated - using useEffect to avoid render issues
    useEffect(() => {
        if (isAuthenticated && admin) {
            // Redirect based on role
            switch (admin.role) {
                case "master_admin":
                    router.push("/master-admin");
                    break;
                case "operation_admin":
                    router.push("/operation-admin");
                    break;
                case "reporting_admin":
                    router.push("/reporting-admin");
                    break;
                default:
                    router.push("/operation-admin");
            }
        }
    }, [isAuthenticated, admin, router]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        e.stopPropagation(); // Add propagation stop
        
        if (!email || !password) {
            toast.error("Email dan password wajib diisi");
            return;
        }

        setIsLoading(true);

        try {
            console.log("Attempting login with", email);
            await login(email, password);
            console.log("Login success, admin data set");
            toast.success("Login berhasil!");
            
            // Trigger manual redirection as well for reliability
            const adminData = localStorage.getItem('access_token');
            if (adminData) {
                // If useEffect doesn't catch it quickly enough
                const response = await api.get("/auth/me");
                const role = response.data.data.role;
                console.log("Manual redirect to", role);
                window.location.href = `/${role.replace('_', '-')}`;
            }
        } catch (error: any) {
            console.error("Login error:", error);
            toast.error(error.message || "Login gagal");
        } finally {
            setIsLoading(false);
        }
    };

    // If authenticated, show a loading state while redirecting
    if (isAuthenticated && !isLoading) {
        return (
            <div className="min-h-screen bg-[#e8e4e0] flex items-center justify-center p-4">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-[#e67e22]"></div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-[#e8e4e0] flex items-center justify-center p-4 md:p-8">
            {/* Main Card Container */}
            <div className="w-full max-w-[1000px] bg-white rounded-3xl shadow-xl overflow-hidden grid grid-cols-1 lg:grid-cols-2 min-h-[600px]">
                
                {/* Left Side - Gradient Panel */}
                <div className="relative flex flex-col justify-between p-8 md:p-10 overflow-hidden">
                    {/* Gradient Background */}
                    <div 
                        className="absolute inset-0"
                        style={{
                            background: "linear-gradient(135deg, #f5e6d8 0%, #f0d4c0 20%, #e8b89a 40%, #e0a080 55%, #d4956e 65%, #f0c8a8 80%, #f5e0d0 100%)",
                        }}
                    />
                    {/* Warm Glow Overlay */}
                    <div 
                        className="absolute top-1/3 left-1/2 -translate-x-1/2 -translate-y-1/3 w-[300px] h-[300px] rounded-full opacity-60"
                        style={{
                            background: "radial-gradient(circle, #e8946a 0%, #f0b890 40%, transparent 70%)",
                        }}
                    />
                    
                    {/* Logo */}
                    <div className="relative z-10">
                        <Image
                            src="/cakli-logo.svg"
                            alt="CakLi Logo"
                            width={120}
                            height={52}
                            priority
                        />
                    </div>

                    {/* Bottom Text */}
                    <div className="relative z-10 mt-auto">
                        <p className="text-sm text-slate-600/80 mb-2">Panel Admin</p>
                        <h2 className="text-2xl md:text-3xl font-bold text-slate-900 leading-tight">
                            Kelola seluruh
                            <br />
                            operasional logistik
                            <br />
                            dalam satu platform.
                        </h2>
                    </div>
                </div>

                {/* Right Side - Login Form */}
                <div className="flex flex-col justify-center p-8 md:p-12 lg:p-14">
                    {/* Decorative Star */}
                    <div className="mb-6">
                        <svg width="28" height="28" viewBox="0 0 28 28" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path d="M14 0L16.5 11.5L28 14L16.5 16.5L14 28L11.5 16.5L0 14L11.5 11.5L14 0Z" fill="#e67e22" />
                        </svg>
                    </div>

                    {/* Heading */}
                    <div className="mb-6">
                        <h1 className="text-2xl md:text-3xl font-bold text-slate-900 tracking-tight">
                            Masuk ke Admin Panel
                        </h1>
                        <p className="text-sm text-slate-500 mt-2 leading-relaxed">
                            Masukkan kredensial Anda untuk mengakses
                            <br />
                            dashboard administrasi CakLi.
                        </p>
                    </div>

                    {/* Form */}
                    <form onSubmit={handleSubmit} method="POST" className="space-y-5">
                        {/* Email Field */}
                        <div className="space-y-2">
                            <Label htmlFor="email" className="text-sm font-medium text-slate-700">
                                Email
                            </Label>
                            <Input
                                id="email"
                                type="email"
                                placeholder="admin@cakli.id"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                disabled={isLoading}
                                className="h-11 border-slate-300 focus:border-[#e67e22] focus:ring-[#e67e22]"
                            />
                        </div>

                        {/* Password Field */}
                        <div className="space-y-2">
                            <Label htmlFor="password" className="text-sm font-medium text-slate-700">
                                Password
                            </Label>
                            <div className="relative">
                                <Input
                                    id="password"
                                    type={showPassword ? "text" : "password"}
                                    placeholder="••••••••"
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                    disabled={isLoading}
                                    className="h-11 pr-10 border-slate-300 focus:border-[#e67e22] focus:ring-[#e67e22]"
                                />
                                <button
                                    type="button"
                                    onClick={() => setShowPassword(!showPassword)}
                                    className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
                                >
                                    {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                                </button>
                            </div>
                        </div>

                        {/* Submit Button */}
                        <Button
                            type="submit"
                            disabled={isLoading}
                            className="w-full h-11 bg-[#e67e22] hover:bg-[#d35400] text-white font-medium rounded-lg transition-colors"
                        >
                            {isLoading ? "Memproses..." : "Masuk"}
                        </Button>
                    </form>

                    {/* Footer Note */}
                    <p className="text-xs text-slate-400 mt-6 text-center">
                        Dengan masuk, Anda menyetujui kebijakan privasi dan syarat layanan CakLi.
                    </p>
                </div>
            </div>
        </div>
    );
}