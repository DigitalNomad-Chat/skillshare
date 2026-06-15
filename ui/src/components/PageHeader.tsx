import { Link } from 'react-router-dom';
import { ArrowLeft, CircleHelp } from 'lucide-react';
import { radius, shadows } from '../design';
import { useT } from '../i18n';
import Tooltip from './Tooltip';

interface PageHeaderProps {
  title: string;
  subtitle?: React.ReactNode;
  icon: React.ReactNode;
  actions?: React.ReactNode;
  className?: string;
  /** Show a styled back button linking to this path */
  backTo?: string;
  /** Help tooltip shown next to the title */
  help?: React.ReactNode;
}

export default function PageHeader({ title, subtitle, icon, actions, className = '', backTo, help }: PageHeaderProps) {
  const t = useT();
  const heading = (
    <div className="flex items-center gap-3">
      {backTo && (
        <Link
          to={backTo}
          className="inline-flex items-center justify-center shrink-0 w-9 h-9 border-2 border-transparent hover:border-muted-dark text-pencil-light hover:text-pencil bg-surface transition-all duration-150 active:scale-95"
          aria-label={t('common.back')}
          style={{ borderRadius: radius.sm, boxShadow: shadows.sm }}
        >
          <ArrowLeft size={18} strokeWidth={2.5} />
        </Link>
      )}
      <div>
        <h2 className="text-2xl md:text-3xl font-bold text-pencil flex items-center gap-2">
          {icon}
          {title}
          {help && (
            <Tooltip content={help} side="top">
              <span className="inline-flex items-center justify-center text-pencil-light hover:text-pencil transition-colors cursor-help ml-1" aria-label={t('common.help')}>
                <CircleHelp size={20} strokeWidth={2} />
              </span>
            </Tooltip>
          )}
        </h2>
        {subtitle && <p className="text-pencil-light mt-1">{subtitle}</p>}
      </div>
    </div>
  );

  return (
    <div
      className={`mb-6 ${actions ? 'flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4' : ''} ${className}`.trim()}
    >
      {heading}
      {actions && <div className="flex items-center gap-2">{actions}</div>}
    </div>
  );
}
