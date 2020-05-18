import React from 'react';
import {Card} from "@material-ui/core";
import { makeStyles } from '@material-ui/core/styles';

function App(props) {
    const calcNumber = () => {
      return 5 + props.foo
    }

    const useStyles = makeStyles((theme) => ({
        card: {
            width: '100px',
            height: '100px',
        },
        container: {
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            padding: '50px',
        },
    }));

    const classes = useStyles();

    return (
        <div className={classes.container}>
            my number is x {calcNumber()}
            <Card classNane={classes.card}>
                test card
            </Card>
        </div>
    );
}

export default App;
