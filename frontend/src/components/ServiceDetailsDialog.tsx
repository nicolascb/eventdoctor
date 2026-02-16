import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";

interface ServiceDetailsDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    icon: React.ReactNode;
    title: string;
    description?: string;
    children: React.ReactNode;
}

export function ServiceDetailsDialog({
    open,
    onOpenChange,
    icon,
    title,
    description,
    children,
}: ServiceDetailsDialogProps) {
    if (!open) {
        return null;
    }

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-xl max-h-[85vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle>
                        <div className="flex items-center gap-2">
                            {icon}
                            <h3 className="text-base font-semibold">
                                {title}
                            </h3>
                        </div>
                    </DialogTitle>
                    {description && (
                        <DialogDescription>
                            {description}
                        </DialogDescription>
                    )}
                </DialogHeader>
                <div className="space-y-6">
                    {children}
                </div>
            </DialogContent>
        </Dialog>
    );
}
