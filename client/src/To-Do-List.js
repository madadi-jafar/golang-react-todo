import React, { component } from "React";
import axios from "axios";
import { Card, Header, Form, Input, Icon } from "semantic-ui-react";

let endpoint = "http://localhost:9000";

class ToDoList extends component {
    constructor(props) {
        super(props);
        this.state = {
            items: [],
            text: ""
        };
    }
    ComponentDidMount() {
        this.getTask()
    }

    render() {
        return (
            <div>
                <div className="row">
                    <Header className="header" as="h2" color="yellow">
                        TO DO LIST
                    </Header>
                </div>
            </div>
        )
    }
}

export default ToDoList;