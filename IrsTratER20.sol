pragma solidity ^0.4.18;

contract IrsTratER20 {

    uint256 constant private MAX_UINT256 = 2**256 - 1;
    mapping (address => uint256) public balances;
    mapping (address => mapping (address => uint256)) public permitted;

    uint256 public totalSupply;             //total existente
    string public name;                   //nombre completo ej:moneda Reports IRSTRAT
    uint8 public decimals;                //cantidad decimales a mostrar.
    string public symbol = "IRT";                 //El simbolo del token: eg IRT

    function IrsTratER20(
        uint256 amount,
        string tn,
        uint8 du,
        string ts
    ) public {
        balances[msg.sender] = amount;               // Asignar al creador todos los tokens
        totalSupply = amount;                        // Registrar este monto inicial
        name = tn;                                   // registrar nombre completo
        decimals = du;                            // registrar los lugares decimales
        symbol = ts;                               // registrar el simbolo inicial
    }

    function transfer(address _to, uint256 _value) public returns (bool success) {
        require(balances[msg.sender] >= _value);
        balances[msg.sender] -= _value;
        balances[_to] += _value;
        Transfer(msg.sender, _to, _value);
        return true;
    }

    function transferFrom(address _from, address _to, uint256 _value) public returns (bool success) {
        uint256 allowance = permitted[_from][msg.sender];
        require(balances[_from] >= _value && allowance >= _value);
        balances[_to] += _value;
        balances[_from] -= _value;
        if (allowance < MAX_UINT256) {
            permitted[_from][msg.sender] -= _value;
        }
        Transfer(_from, _to, _value);
        return true;
    }

    function balanceOf(address _owner) public view returns (uint256 balance) {
        return balances[_owner];
    }

    function approve(address _spender, uint256 _value) public returns (bool success) {
        permitted[msg.sender][_spender] = _value;
        Approval(msg.sender, _spender, _value);
        return true;
    }

    function allowance(address _owner, address _spender) public view returns (uint256 remaining) {
        return permitted[_owner][_spender];
    }

    event Transfer(address indexed _from, address indexed _to, uint256 _value);
    event Approval(address indexed _owner, address indexed _spender, uint256 _value);
}
