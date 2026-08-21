import { Link } from "react-router-dom";
import { ArrowRight, Phone, Clock } from "lucide-react";

export function CtaDark() {
  return (
    <section className="bg-[#111827] dark:bg-[#0B1120] border-y border-[#1F2937]">
      <div className="max-w-[1280px] mx-auto px-6 py-24">
        <div className="max-w-2xl">
          <div className="accent-line mb-6" />
          <h2 className="text-section text-white mb-4">
            Ready to join<br />the network?
          </h2>
          <p className="text-body text-gray-400 mb-10 max-w-lg">
            Connect your healthcare organization to the Curexal network. Onboarding takes days, not months, and every organization retains full ownership of its data.
          </p>
          <div className="flex flex-col sm:flex-row items-start gap-3">
            <Link to="/book-demo">
              <button
                className="flex items-center gap-2 px-6 py-3 rounded-[10px] bg-[#0F766E] hover:bg-[#14B8A6] text-white text-sm font-semibold transition-colors cursor-pointer border-0"
                style={{ boxShadow: "0 1px 2px rgba(0,0,0,0.2), 0 4px 16px rgba(15,118,110,0.25)" }}
              >
                <span>Book Demo</span>
                <span className="text-[10px] bg-amber-500/20 text-amber-300 font-extrabold px-1.5 py-0.5 rounded flex items-center gap-1">
                  <Clock className="w-3 h-3" />
                  Soon
                </span>
                <ArrowRight className="h-4 w-4 ml-1" />
              </button>
            </Link>
            <a href="mailto:sales@curexal.com">
              <button className="flex items-center gap-2 px-6 py-3 rounded-[10px] border border-[#374151] text-gray-300 text-sm font-semibold hover:border-[#4B5563] hover:bg-[#1F2937] transition-colors cursor-pointer bg-transparent">
                <Phone className="h-4 w-4" />
                Contact Sales
              </button>
            </a>
          </div>
        </div>
      </div>
    </section>
  );
}
