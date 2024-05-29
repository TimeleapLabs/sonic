package sfclib

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// GetContractBin is SFCLib contract genesis implementation bin code
// Has to be compiled with flag bin-runtime
// Built from opera-sfc 424031c81a77196f4e9d60c7d876032dd47208ce, solc 0.5.17+commit.d19bba13.Emscripten.clang, optimize-runs 200
func GetContractBin() []byte {
	return hexutil.MustDecode("0x608060405234801561001057600080fd5b50600436106101e55760003560e01c8063893675c61161010f578063c5f956af116100a2578063cfd4766311610071578063cfd4766314610613578063cfdbb7cd1461064c578063d96ed50514610685578063f2fde38b1461068d576101e5565b8063c5f956af146105c1578063c65ee0e1146105c9578063c7be95de146105e6578063cc8343aa146105ee576101e5565b806396c7ee46116100de57806396c7ee4614610485578063a86a056f146104e4578063b5d896271461051d578063b810e41114610588576101e5565b8063893675c6146104515780638b0e9f3f146104595780638da5cb5b146104615780638f32d59b14610469576101e5565b806339b80c0011610187578063766718081161015657806376671808146103925780637ad797871461039a5780637cacb1d6146103b7578063854873e1146103bf576101e5565b806339b80c00146102f45780635fab23a814610349578063670322f814610351578063715018a61461038a576101e5565b806318160ddd116101c357806318160ddd1461027f5780631f2701521461028757806328f73148146102e457806334d722c9146102ec576101e5565b80630135b1db146101ea5780630553fd5b1461022f5780630e559d821461024e575b600080fd5b61021d6004803603602081101561020057600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166106c0565b60408051918252519081900360200190f35b61024c6004803603602081101561024557600080fd5b50356106d2565b005b610256610764565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61021d610780565b6102c66004803603606081101561029d57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060208101359060400135610786565b60408051938452602084019290925282820152519081900360600190f35b61021d6107b8565b6102566107be565b6103116004803603602081101561030a57600080fd5b50356107da565b604080519788526020880196909652868601949094526060860192909252608085015260a084015260c0830152519081900360e00190f35b61021d61081c565b61021d6004803603604081101561036757600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610822565b61024c610873565b61021d610955565b61024c600480360360208110156103b057600080fd5b503561095e565b61021d6109ed565b6103dc600480360360208110156103d557600080fd5b50356109f3565b6040805160208082528351818301528351919283929083019185019080838360005b838110156104165781810151838201526020016103fe565b50505050905090810190601f1680156104435780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610256610aac565b61021d610ac8565b610256610ace565b610471610aea565b604080519115158252519081900360200190f35b6104be6004803603604081101561049b57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610b08565b604080519485526020850193909352838301919091526060830152519081900360800190f35b61021d600480360360408110156104fa57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610b3a565b61053a6004803603602081101561053357600080fd5b5035610b57565b604080519788526020880196909652868601949094526060860192909252608085015260a084015273ffffffffffffffffffffffffffffffffffffffff1660c0830152519081900360e00190f35b6102c66004803603604081101561059e57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610baa565b610256610bd6565b61021d600480360360208110156105df57600080fd5b5035610bf2565b61021d610c04565b61024c6004803603604081101561060457600080fd5b50803590602001351515610c0a565b61021d6004803603604081101561062957600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610e77565b6104716004803603604081101561066257600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610e94565b61021d610f52565b61024c600480360360208110156106a357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16610f58565b60696020526000908152604090205481565b60775473ffffffffffffffffffffffffffffffffffffffff16331461075857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f63616c6c6572206973206e6f742061206d696e74657200000000000000000000604482015290519081900360640190fd5b61076181610fd4565b50565b607c5473ffffffffffffffffffffffffffffffffffffffff1681565b60765481565b607160209081526000938452604080852082529284528284209052825290208054600182015460029092015490919083565b606d5481565b60775473ffffffffffffffffffffffffffffffffffffffff1681565b607860205280600052604060002060009150905080600701549080600801549080600901549080600a01549080600b01549080600c01549080600d0154905087565b606e5481565b600061082e8383610e94565b61083a5750600061086d565b5073ffffffffffffffffffffffffffffffffffffffff821660009081526073602090815260408083208484529091529020545b92915050565b61087b610aea565b6108e657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b60335460405160009173ffffffffffffffffffffffffffffffffffffffff16907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3603380547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055565b60675460010190565b60775473ffffffffffffffffffffffffffffffffffffffff1633146109e457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f63616c6c6572206973206e6f742061206d696e74657200000000000000000000604482015290519081900360640190fd5b61076181610fed565b60675481565b606a6020908152600091825260409182902080548351601f60027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff61010060018616150201909316929092049182018490048402810184019094528084529091830182828015610aa45780601f10610a7957610100808354040283529160200191610aa4565b820191906000526020600020905b815481529060010190602001808311610a8757829003601f168201915b505050505081565b60835473ffffffffffffffffffffffffffffffffffffffff1681565b606c5481565b60335473ffffffffffffffffffffffffffffffffffffffff1690565b60335473ffffffffffffffffffffffffffffffffffffffff16331490565b607360209081526000928352604080842090915290825290208054600182015460028301546003909301549192909184565b607060209081526000928352604080842090915290825290205481565b6068602052600090815260409020805460018201546002830154600384015460048501546005860154600690960154949593949293919290919073ffffffffffffffffffffffffffffffffffffffff1687565b607460209081526000928352604080842090915290825290208054600182015460029092015490919083565b60805473ffffffffffffffffffffffffffffffffffffffff1681565b607b6020526000908152604090205481565b606b5481565b610c1382611000565b610c7e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f76616c696461746f7220646f65736e2774206578697374000000000000000000604482015290519081900360640190fd5b60008281526068602052604090206003810154905415610c9c575060005b606654604080517fa4066fbe0000000000000000000000000000000000000000000000000000000081526004810186905260248101849052905173ffffffffffffffffffffffffffffffffffffffff9092169163a4066fbe9160448082019260009290919082900301818387803b158015610d1657600080fd5b505af1158015610d2a573d6000803e3d6000fd5b50505050818015610d3a57508015155b15610e72576066546000848152606a60205260409081902081517f242a6e3f0000000000000000000000000000000000000000000000000000000081526004810187815260248201938452825460027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60018316156101000201909116046044830181905273ffffffffffffffffffffffffffffffffffffffff9095169463242a6e3f94899493909160649091019084908015610e385780601f10610e0d57610100808354040283529160200191610e38565b820191906000526020600020905b815481529060010190602001808311610e1b57829003601f168201915b50509350505050600060405180830381600087803b158015610e5957600080fd5b505af1158015610e6d573d6000803e3d6000fd5b505050505b505050565b607260209081526000928352604080842090915290825290205481565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260736020908152604080832084845290915281206002015415801590610f05575073ffffffffffffffffffffffffffffffffffffffff8316600090815260736020908152604080832085845290915290205415155b8015610f4b575073ffffffffffffffffffffffffffffffffffffffff83166000908152607360209081526040808320858452909152902060020154610f48611017565b11155b9392505050565b607f5481565b610f60610aea565b610fcb57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b6107618161101b565b607654610fe7908263ffffffff61111516565b60765550565b607654610fe7908263ffffffff61118916565b600090815260686020526040902060050154151590565b4290565b73ffffffffffffffffffffffffffffffffffffffff8116611087576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260268152602001806112786026913960400191505060405180910390fd5b60335460405173ffffffffffffffffffffffffffffffffffffffff8084169216907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3603380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b600082820183811015610f4b57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b6000610f4b83836040518060400160405280601e81526020017f536166654d6174683a207375627472616374696f6e206f766572666c6f7700008152506000818484111561126f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b8381101561123457818101518382015260200161121c565b50505050905090810190601f1680156112615780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b50505090039056fe4f776e61626c653a206e6577206f776e657220697320746865207a65726f2061646472657373a265627a7a72315820a4e817087310a426b106e703d6cd1a56132957d6108c6cb1f854fff562551cd264736f6c63430005110032")
}

// ContractAddress is the SFCLib contract address
var ContractAddress = common.HexToAddress("0xfc01face00000000000000000000000000000000")
