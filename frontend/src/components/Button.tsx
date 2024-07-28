interface ButtonProps extends React.ComponentProps<'button'> {
  loading: boolean
}

export const Button = ({ content, loading, ...props }: ButtonProps) => (
  <button
    className="btn btn-primary disabled:bg-gray-700 disabled:text-gray-400"
    {...props}
  >
    {loading ? <span className="loading loading-spinner"></span> : content}
  </button>
)
