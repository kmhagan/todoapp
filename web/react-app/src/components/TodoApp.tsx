import React, { FunctionComponent, useState } from 'react'

interface Props {
}

const TodoApp: FunctionComponent<Props> = (props) => {
  const [data, setData] = useState(Array<string>())
  const [current, setCurrent] = useState("")

  const handleChange = function(event: React.ChangeEvent<HTMLInputElement>) {
    setCurrent(event.target.value)
  }

  const handleSubmit = function(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setData([...data, current])
    setCurrent("")
  }

  return (
    <div className="todoListMain">
      <div className="header">
        <form onSubmit={handleSubmit}>
          <input placeholder="Task" value={current} onChange={handleChange} />
          <button type="submit"> Add Task </button>
          {data.map((value, i) => (
            <li key={i}>
              {value}
            </li>
          ))}
        </form>
      </div>
    </div>
  )

}

export default TodoApp;
