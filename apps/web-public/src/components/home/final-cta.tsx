import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Activity, ArrowRight, Clock } from "lucide-react";

export function FinalCta() {
  return (
    <section id="about" className="py-24 px-6 border-t border-border">
      <div className="max-w-4xl mx-auto text-center">
        <div className="relative rounded-3xl border border-primary/20 bg-gradient-to-b from-primary/10 to-primary/2 p-16 overflow-hidden">
          <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
            <div className="w-[500px] h-[300px] bg-primary/15 blur-[100px] rounded-full" />
          </div>
          <div className="relative z-10">
            <Activity className="h-12 w-12 text-primary mx-auto mb-6 animate-glow-pulse" />
            <h2 className="text-4xl font-black mb-4 text-foreground">
              Deploy Curexal to your{" "}
              <span className="brand-gradient">diagnostics laboratory</span>
            </h2>
            <p className="text-muted-foreground text-lg mb-10 max-w-xl mx-auto">
              Begin managing your laboratory workspace, inviting clinical pathology teams, configuring test profiles, and maintaining regulatory audit logging today.
            </p>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
              <Link to="/book-demo" id="cta-book-demo">
                <Button
                  size="lg"
                  className="px-8 py-6 text-base bg-[#0F766E] hover:bg-[#115E59] text-white shadow-xl shadow-[#0F766E]/30 gap-2 rounded-xl border-0 cursor-pointer font-bold"
                >
                  <span>Book Demo</span>
                  <span className="text-[10px] bg-amber-500/20 text-amber-300 font-extrabold px-2 py-0.5 rounded flex items-center gap-1">
                    <Clock className="w-3 h-3" />
                    Soon
                  </span>
                  <ArrowRight className="h-4 w-4 ml-1" />
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
