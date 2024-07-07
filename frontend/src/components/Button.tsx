type ButtonProps = React.ComponentProps<'button'> & { loading: boolean };

const Button = ({ onClick, content, loading }: ButtonProps) => (
  <button className="btn btn-primary" onClick={onClick}>
    {loading ? <span className="loading loading-spinner"></span> : content}
  </button>
);

export default Button;
