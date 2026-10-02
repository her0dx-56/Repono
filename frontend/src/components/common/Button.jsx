export default function Button({children,variant='dark',className='',...props}){return <button className={`btn btn-${variant} ${className}`} {...props}>{children}</button>}
