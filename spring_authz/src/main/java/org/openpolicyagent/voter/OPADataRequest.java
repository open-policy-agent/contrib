package org.openpolicyagent.voter;

import java.util.Map;

public class OPADataRequest {

    Map<String, Object> input;

    public OPADataRequest(Map<String, Object> input) {
        this.input = input;
    }

    public Map<String, Object> getInput() {
        return this.input;
    }

}