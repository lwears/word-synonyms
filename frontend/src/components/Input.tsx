import clsx from 'clsx'

import type { UseFormRegister, FieldError } from 'react-hook-form'
import type { FormData } from '../types'

interface InputProps extends React.ComponentProps<'input'> {
  register: UseFormRegister<FormData>
  error: FieldError | undefined
  name: 'synonym' | 'word'
  title: string
}

const Input: React.FC<InputProps> = ({
  type,
  placeholder,
  name,
  register,
  error,
  title,
}) => (
  <>
    <div className="label">
      <span className="label-text text-primary">{title}</span>
    </div>
    <input
      className={clsx(
        'input input-bordered w-full max-w-xs',
        error && 'input-error'
      )}
      type={type}
      placeholder={placeholder}
      {...register(name)}
    />
    {error && (
      <div className="label">
        <span className="label-text-alt text-error">{error.message}</span>
      </div>
    )}
  </>
)
export default Input

{
  /* <label className="form-control w-full max-w-xs">
  <div className="label">
    <span className="label-text">Enter Word</span>
  </div>
  <input
    id="word"
    type="text"
    placeholder="Type word here"
    className={clsx(
      'input input-bordered w-full max-w-xs',
      errors.word && 'input-error'
    )}
    required
    disabled={loading}
    {...register('word')}
  />
  {errors.word && (
    <div className="label">
      <span className="label-text-alt text-error">{errors.word.message}</span>
    </div>
  )}
</label>; */
}
