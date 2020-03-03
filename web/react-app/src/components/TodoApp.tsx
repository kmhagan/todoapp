import React from 'react'

interface TodoAppState {
  data: Array<string>,
  current: string,
}

class TodoApp extends React.Component<{}, TodoAppState> {
  constructor(props: TodoAppState) {
    super(props)
    this.state = {
      data: [],
      current: "",
    }
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event: React.ChangeEvent<HTMLInputElement>) {
    this.setState({current: event.target.value});
  }

  handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    this.setState({
      data: [...this.state.data, this.state.current],
      current: "",
    });
  }

  render() {
    return (
      <div className="todoListMain">
        <div className="header">
          <form onSubmit={this.handleSubmit}>
            <input placeholder="Task" value={this.state.current} onChange={this.handleChange} />
            <button type="submit"> Add Task </button>
            {this.state.data.map((data, i) => (
              <li key={i}>
                {data}
              </li>
            ))}
          </form>
        </div>
      </div>
    )
  }
}

export default TodoApp;