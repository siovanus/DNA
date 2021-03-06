/*
 * Copyright (C) 2018 The DNA Authors
 * This file is part of The DNA library.
 *
 * The DNA is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The DNA is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The DNA.  If not, see <http://www.gnu.org/licenses/>.
 */
package ontid

import (
	"github.com/dnaproject2/DNA/smartcontract/service/native"
	"github.com/dnaproject2/DNA/smartcontract/service/native/utils"
)

func Init() {
	native.Contracts[utils.OntIDContractAddress] = RegisterIDContract
}

func RegisterIDContract(srvc *native.NativeService) {
	srvc.Register("regIDWithPublicKey", regIdWithPublicKey)
	srvc.Register("regIDWithController", regIdWithController)
	srvc.Register("revokeID", revokeID)
	srvc.Register("revokeIDByController", revokeIDByController)
	srvc.Register("removeController", removeController)
	srvc.Register("addRecovery", addRecovery)
	srvc.Register("changeRecovery", changeRecovery)
	srvc.Register("addKey", addKey)
	srvc.Register("removeKey", removeKey)
	srvc.Register("addKeyByController", addKeyByController)
	srvc.Register("removeKeyByController", removeKeyByController)
	srvc.Register("addKeyByRecovery", addKeyByRecovery)
	srvc.Register("removeKeyByRecovery", removeKeyByRecovery)
	srvc.Register("regIDWithAttributes", regIdWithAttributes)
	srvc.Register("addAttributes", addAttributes)
	srvc.Register("removeAttribute", removeAttribute)
	srvc.Register("addAttributesByController", addAttributesByController)
	srvc.Register("removeAttributeByController", removeAttributeByController)
	srvc.Register("verifySignature", verifySignature)
	srvc.Register("verifyController", verifyController)
	srvc.Register("getPublicKeys", GetPublicKeys)
	srvc.Register("getKeyState", GetKeyState)
	srvc.Register("getAttributes", GetAttributes)
	srvc.Register("getDDO", GetDDO)
	return
}
