import React from 'react';
import { withStyles } from '@material-ui/core/styles';
import TextField from "@material-ui/core/TextField";
import Button from "@material-ui/core/Button";
import { WrapperMessage, ServerConnectionSuccessMessage } from "proto/messages_pb.js"

class App extends React.Component {
    state = {username: '', email: ''}

    handleUsernameChange = (event) => {
        this.setState({username: event.target.value})
    }
    handleEmailChange = (event) => {
        this.setState({email: event.target.value})
    }
    handleStartClick = () => {
        console.log(this.state.username, this.state.email)
        const socket = new WebSocket(
            `ws://localhost:8000/ws?email=${this.state.email}&name=${this.state.username}`
        );

        socket.addEventListener('open', function (event) {
            socket.send('Hello Server!');
        });

// Listen for messages
        socket.addEventListener('message', function (event) {
            console.log('Message from server ', event.data);

        });
        console.log('success')
    }

    render() {
        const {classes} = this.props;
        return (
            <div className={classes.container}>
                <div className={classes.header}>Welcome to Bot Or Not</div>
                <div>
                    <TextField
                        id="standard-basic"
                        label="Username"
                        value={this.state.username}
                        onChange={this.handleUsernameChange}
                    />
                </div>
                <div>
                    <TextField
                        id="standard-basic"
                        label="Email"
                        value={this.state.email}
                        onChange={this.handleEmailChange}
                    />
                </div>
                <div>
                    <Button onClick={this.handleStartClick}>
                        Start
                    </Button>
                </div>
            </div>
        );
    }
}
const styles = {
    card: {
        width: '100px',
        height: '100px',
    },
    container: {
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        padding: '50px',
        width: '100%'
    },
    header: {
        fontSize: '25px',
    }
}

export default withStyles(styles)(App);
