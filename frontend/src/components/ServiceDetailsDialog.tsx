import {
    Sheet,
    SheetContent,
    SheetHeader,
    SheetTitle,
    SheetDescription,
} from "@/components/ui/sheet";

interface ServiceDetailsDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    icon: React.ReactNode;
    title: string;
    description?: string;
    children: React.ReactNode;
    className?: string;
}

export function ServiceDetailsDialog({
    open,
    onOpenChange,
    icon,
    title,
    description,
    children,
    className,
}: ServiceDetailsDialogProps) {
    if (!open) {
        return null;
    }

    return (
        <Sheet open={open} onOpenChange={onOpenChange}>
            <SheetContent className={`w-full overflow-y-auto p-0 flex flex-col ${className || "sm:max-w-xl"}`}>
                <div className="p-6 pb-4 border-b border-border bg-muted/20">
                    <SheetHeader>
                        <div className="flex items-center gap-3">
                            <div className="rounded-md bg-primary/10 p-2 text-primary">
                                {icon}
                            </div>
                            <div className="text-left">
                                <SheetTitle className="text-xl font-mono">
                                    {title}
                                </SheetTitle>
                                {description && (
                                    <SheetDescription className="text-xs truncate max-w-[300px]" title={description}>
                                        {description}
                                    </SheetDescription>
                                )}
                            </div>
                        </div>
                    </SheetHeader>
                </div>
                <div className="p-6 flex-1 bg-background space-y-6">
                    {children}
                </div>
            </SheetContent>
        </Sheet>
    );
}
