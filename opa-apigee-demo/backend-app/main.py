from flask import Flask, jsonify

app = Flask(__name__)


# Test claims

claims = {
    '1001':  {
        'enrolleeId': '1001',
        'enrolleePrimaryId': '1001',
        'enrolleeAge': 32,
        'canPrimaryViewClaims': False,
        'enrolleeName': 'Alice Opa',
        'enrolleeType': 'Primary',
        'enrolleeClaimSummaryList': [
            {
                'claimId': 'Alice - 20172376088994',
                'claimStatus': 'Claim Paid',
                'claimStatusCode': 'Y',
                'dateOfService': '08/25/2017',
                'provider': 'Zara Medical Center',
                'enrolleeResponsibilityAmount': 100,
                'procedure': 'Cataract surgery',
            },
            {
                'claimId': 'Alice - 20172376088995',
                'claimStatus': 'Claim Paid',
                'claimStatusCode': 'Y',
                'dateOfService': '10/25/2017',
                'provider': 'Louis Vuitton Medical Center',
                'enrolleeResponsibilityAmount': 75,
                'procedure': 'Appendectomy',
            },
            {
                'claimId': 'Alice - 20172376088996',
                'claimStatus': 'Claim Not Paid',
                'claimStatusCode': 'N',
                'dateOfService': '04/25/2020',
                'provider': 'Chanel Medical Center',
                'enrolleeResponsibilityAmount': 1200,
                'procedure': 'Low back pain surgery',
            }
        ]
    },
    '1002': {
        'enrolleeId': '1002',
        'enrolleePrimaryId': '1001',
        'enrolleeAge': 50,
        'canPrimaryViewClaims': True,
        'enrolleeName': 'Bob Opa',
        'enrolleeType': 'Secondary',
        'enrolleeClaimSummaryList': [
            {
                'claimId': 'Bob - 20172376088995',
                'claimStatus': 'Claim Paid',
                'claimStatusCode': 'Y',
                'dateOfService': '04/25/2019',
                'provider': 'Zara Medical Center',
                'enrolleeResponsibilityAmount': 0,
                'procedure': 'Coronary artery bypass',
            }
        ]
    },
    '1003': {
        'enrolleeId': '1003',
        'enrolleePrimaryId': '1001',
        'enrolleeAge': 15,
        'canPrimaryViewClaims': False,
        'enrolleeName': 'John Opa',
        'enrolleeType': 'Secondary',
        'enrolleeClaimSummaryList': [
            {
                'claimId': 'John - 20172376089990',
                'claimStatus': 'Claim Paid',
                'claimStatusCode': 'Y',
                'dateOfService': '04/14/2019',
                'provider': 'Gucci Medical Center',
                'enrolleeResponsibilityAmount': 0,
                'procedure': 'CT Scan',
            },
            {
                'claimId': 'John - 20172376089991',
                'claimStatus': 'Claim Paid',
                'claimStatusCode': 'Y',
                'dateOfService': '04/14/2020',
                'provider': 'Dolce Medical Center',
                'enrolleeResponsibilityAmount': 0,
                'procedure': 'MRI',
            }
        ]
    },
    '1004': {
        'enrolleeId': '1004',
        'enrolleePrimaryId': '1001',
        'enrolleeAge': 21,
        'canPrimaryViewClaims': False,
        'enrolleeName': 'Jane Opa',
        'enrolleeType': 'Secondary',
        'enrolleeClaimSummaryList': [
            {
                'claimId': 'Jane - 20172376089998',
                'claimStatus': 'Claim Paid',
                'claimStatusCode': 'Y',
                'dateOfService': '03/15/2019',
                'provider': 'Zara Medical Center',
                'enrolleeResponsibilityAmount': 200,
                'procedure': 'Low back pain surgery',
            }
        ]
    }
}


@app.route('/', methods=['GET'])
def homepage():
    return "<h1>Claims Service</h1><p>This service holds claims of Acme Health Care members.</p>"


@app.route('/v1/claims/enrollee/<enrolleeID>', methods=['GET'])
def get_enrollee_claims(enrolleeID):
    if enrolleeID not in claims:
        return 'enrolleeID {} does not exist'.format(enrolleeID), 404
    return jsonify(claims[enrolleeID])


if __name__ == "__main__":
    app.run()
