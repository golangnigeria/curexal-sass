import { type ReactNode } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { authClient } from "@/lib/auth-client";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { useTheme } from "@/components/theme-provider";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  LayoutDashboard,
  Users,
  Shield,
  CreditCard,
  Settings,
  Menu,
  Building2,
  Activity,
  Bell,
  Search,
  LogOut,
  ChevronRight,
  Sun,
  Moon,
} from "lucide-react";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";

const navigation = [
  {
    group: "Workspace",
    items: [
      { name: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
      { name: "Members", href: "/members", icon: Users },
      { name: "Roles & Permissions", href: "/roles", icon: Shield },
      { name: "Settings", href: "/settings", icon: Settings },
    ],
  },
];

interface AppLayoutProps {
  children: ReactNode;
}

function NavItem({ name, href, icon: Icon }: { name: string; href: string; icon: any }) {
  const location = useLocation();
  const isActive = location.pathname === href;

  return (
    <Link to={href}>
      <div
        className={cn(
          "group flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all duration-150 cursor-pointer border border-transparent",
          isActive
            ? "bg-primary/10 text-primary border-primary/20 dark:bg-primary/15"
            : "text-muted-foreground hover:text-foreground hover:bg-accent/50"
        )}
      >
        <Icon
          className={cn(
            "h-4 w-4 flex-shrink-0 transition-colors",
            isActive ? "text-primary" : "text-muted-foreground group-hover:text-foreground"
          )}
        />
        {name}
        {isActive && (
          <ChevronRight className="h-3.5 w-3.5 ml-auto text-primary/60" />
        )}
      </div>
    </Link>
  );
}

function SidebarContent({ session }: { session: any }) {
  return (
    <div className="flex flex-col h-full py-6 px-4">
      {/* Logo */}
      <Link to="/" className="flex items-center gap-2.5 px-2 mb-8 group">
        <div className="w-8 h-8 rounded-xl bg-primary/20 flex items-center justify-center group-hover:bg-primary/30 transition-colors">
          <Activity className="h-4 w-4 text-primary" />
        </div>
        <div>
          <span className="font-bold text-base brand-gradient">Curexal</span>
          <p className="text-[10px] text-muted-foreground font-medium -mt-0.5">Control Plane</p>
        </div>
      </Link>

      {/* Navigation groups */}
      <nav className="flex-1 space-y-6 overflow-y-auto">
        {navigation.map((group) => (
          <div key={group.group}>
            <p className="px-3 mb-2 text-[10px] font-semibold text-muted-foreground/60 uppercase tracking-widest">
              {group.group}
            </p>
            <div className="space-y-0.5">
              {group.items.map((item) => (
                <NavItem key={item.href} {...item} />
              ))}
            </div>
          </div>
        ))}
      </nav>

      {/* Bottom: status pill */}
      <div className="pt-4 border-t border-border">
        <div className="flex items-center gap-2 px-3 py-2 rounded-xl bg-emerald-500/8 border border-emerald-500/15">
          <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
          <span className="text-xs text-emerald-400 font-medium">All systems normal</span>
        </div>
      </div>
    </div>
  );
}

export function AppLayout({ children }: AppLayoutProps) {
  const navigate = useNavigate();
  const { data: session } = authClient.useSession();
  const { theme, setTheme } = useTheme();

  const handleSignOut = async () => {
    await authClient.signOut();
    navigate("/");
  };

  const userInitials = session?.user?.name
    ?.split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase()
    .slice(0, 2) ?? "?";

  return (
    <div className="min-h-screen bg-background flex text-foreground">
      {/* ── Desktop Sidebar ───────────────────── */}
      <aside className="hidden md:flex md:w-60 md:flex-col fixed inset-y-0 border-r border-border bg-sidebar text-sidebar-foreground">
        <SidebarContent session={session} />
      </aside>

      {/* ── Main Content ─────────────────────── */}
      <div className="flex-1 md:pl-60 min-w-0">
        {/* Top Bar */}
        <header className="sticky top-0 z-40 flex h-14 items-center gap-4 border-b border-border bg-background/80 backdrop-blur-xl px-6">
          {/* Mobile burger */}
          <Sheet>
            <SheetTrigger asChild>
              <Button variant="ghost" size="icon" className="md:hidden text-muted-foreground hover:text-foreground">
                <Menu className="h-5 w-5" />
              </Button>
            </SheetTrigger>
            <SheetContent side="left" className="w-60 p-0 bg-sidebar border-r border-border text-sidebar-foreground">
              <SidebarContent session={session} />
            </SheetContent>
          </Sheet>

          {/* Search placeholder */}
          <div className="hidden sm:flex flex-1 max-w-sm">
            <div className="flex items-center gap-2 w-full px-3 py-1.5 rounded-lg border border-border bg-muted/50 text-sm text-muted-foreground cursor-pointer hover:border-border/80 transition-colors">
              <Search className="h-3.5 w-3.5" />
              <span>Search...</span>
              <kbd className="ml-auto text-[10px] border border-border rounded px-1.5 py-0.5 text-muted-foreground">⌘K</kbd>
            </div>
          </div>

          <div className="flex-1" />

          {/* Right actions */}
          <div className="flex items-center gap-3">
            {/* Notifications */}
            <Button variant="ghost" size="icon" className="h-8 w-8 text-slate-500 hover:text-slate-900 dark:hover:text-white relative">
              <Bell className="h-4 w-4" />
              <span className="absolute top-1 right-1 w-1.5 h-1.5 rounded-full bg-primary" />
            </Button>

            {/* Theme Toggle */}
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground hover:text-foreground relative">
                  <Sun className="h-4 w-4 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
                  <Moon className="absolute h-4 w-4 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
                  <span className="sr-only">Toggle theme</span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="bg-popover border-border">
                <DropdownMenuItem
                  className={cn(
                    "cursor-pointer text-muted-foreground hover:text-foreground focus:bg-accent/50",
                    theme === "light" && "text-primary font-medium"
                  )}
                  onClick={() => setTheme("light")}
                >
                  <Sun className="h-4 w-4 mr-2" />
                  Light
                </DropdownMenuItem>
                <DropdownMenuItem
                  className={cn(
                    "cursor-pointer text-muted-foreground hover:text-foreground focus:bg-accent/50",
                    theme === "dark" && "text-primary font-medium"
                  )}
                  onClick={() => setTheme("dark")}
                >
                  <Moon className="h-4 w-4 mr-2" />
                  Dark
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>

            {/* User menu */}
            {session?.user && (
              <DropdownMenu>
                <DropdownMenuTrigger className="focus:outline-none" id="user-menu-trigger">
                  <Avatar className="h-8 w-8 cursor-pointer ring-2 ring-transparent hover:ring-primary/30 transition-all">
                    <AvatarImage src={session.user.image || undefined} alt={session.user.name} />
                    <AvatarFallback className="bg-primary/20 text-primary text-xs font-bold">
                      {userInitials}
                    </AvatarFallback>
                  </Avatar>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-56 bg-popover border-border">
                  <DropdownMenuLabel className="font-normal">
                    <div className="flex flex-col gap-0.5">
                      <p className="text-sm font-semibold text-foreground">{session.user.name}</p>
                      <p className="text-xs text-muted-foreground">{session.user.email}</p>
                    </div>
                  </DropdownMenuLabel>
                  <DropdownMenuSeparator className="bg-border" />
                  <DropdownMenuItem
                    className="cursor-pointer text-muted-foreground hover:text-foreground focus:bg-accent/50"
                    onClick={() => navigate("/settings")}
                  >
                    <Settings className="h-3.5 w-3.5 mr-2" />
                    Settings
                  </DropdownMenuItem>
                  <DropdownMenuSeparator className="bg-border" />
                  <DropdownMenuItem
                    className="cursor-pointer text-red-500 dark:text-red-400 focus:bg-red-500/8 focus:text-red-600"
                    onClick={handleSignOut}
                    id="sign-out-btn"
                  >
                    <LogOut className="h-3.5 w-3.5 mr-2" />
                    Sign out
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            )}
          </div>
        </header>

        {/* Page Content */}
        <main className="py-8 px-6">
          {children}
        </main>
      </div>
    </div>
  );
}