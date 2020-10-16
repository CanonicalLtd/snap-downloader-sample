import React, {Component} from 'react';
import {formatError, parseRoute} from "./components/Utils";
import Header from "./components/Header";
import Home from "./components/Home";
import Footer from "./components/Footer";
import api from "./components/api";


class App extends Component  {
    constructor(props) {
        super(props)
        this.state = {
            macaroon: {},
        }
    }

    componentDidMount() {
        this.getAuth();
    }

    getAuth() {
        api.authGet().then(response => {
            this.setState({macaroon: response.data.record})
        })
        .catch(e => {
            console.log(formatError(e.response.data))
            this.setState({error: formatError(e.response.data), message: ''});
        })
    }

    handleLogin = () => {
        this.getAuth();
    }

    render() {
        const r = parseRoute()

        return (
            <div>
              <Header />

              {r.section===''? <Home macaroon={this.state.macaroon} onLogin={this.handleLogin()} /> : ''}
              {/*{r.section==='builds'? <BuildLog buildId={r.sectionId} /> : ''}*/}
              {/*{r.section==='settings'? <Settings /> : ''}*/}

              <Footer />
            </div>
        );
    }
}

export default App;
