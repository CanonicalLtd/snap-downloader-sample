import React, {Component} from 'react';
import {formatError, parseRoute} from "./components/Utils";
import Header from "./components/Header";
import Home from "./components/Home";
import Footer from "./components/Footer";
import api from "./components/api";
import Login from "./components/Login";
import Settings from "./components/Settings";


class App extends Component  {
    constructor(props) {
        super(props)
        this.state = {
            macaroon: null,
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
            this.setState({error: formatError(e.response.data), message: '', macaroon: null});
        })
    }

    postLogin = () => {
        this.getAuth();
    }

    renderHome() {
        if (!this.state.macaroon) {
            return <Login onLogin={this.postLogin} />
        }
        if (!this.state.macaroon['Snap-Device-Store']) {
            return <Login onLogin={this.postLogin} />
        }
        return <Home macaroon={this.state.macaroon} />
    }

    renderSettings() {
        if (!this.state.macaroon) {
            return <Login onLogin={this.postLogin} />
        }
        if (!this.state.macaroon['Snap-Device-Store']) {
            return <Login onLogin={this.postLogin} />
        }
        return <Settings />
    }

    render() {
        const r = parseRoute()

        return (
            <div>
              <Header />

              {r.section==='' ? this.renderHome() : ''}
              {r.section==='login'? <Login onLogin={this.postLogin} /> : ''}
              {r.section==='settings' ? this.renderSettings(): ''}

              <Footer />
            </div>
        );
    }
}

export default App;
