import { Input } from '@/components/ui/input';
import { Search, X } from 'lucide-react';

interface SearchInputProps {
    value: string;
    onChange: (value: string) => void;
    placeholder?: string;
    resultCount?: number;
    totalCount?: number;
    className?: string;
}

export function SearchInput({
    value,
    onChange,
    placeholder = 'Search...',
    resultCount,
    totalCount,
    className,
}: SearchInputProps) {
    const showCount = value.trim().length > 0 && resultCount !== undefined && totalCount !== undefined;

    return (
        <div className={`relative ${className || ''}`}>
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
                type="text"
                value={value}
                onChange={(e) => onChange(e.target.value)}
                placeholder={placeholder}
                className="pl-10 pr-10 h-9 bg-card border-border text-sm"
            />
            {value && (
                <button
                    type="button"
                    onClick={() => onChange('')}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
                >
                    <X className="h-4 w-4" />
                </button>
            )}
            {showCount && (
                <p className="absolute -bottom-5 left-0 text-xs text-muted-foreground">
                    {resultCount} of {totalCount} results
                </p>
            )}
        </div>
    );
}
